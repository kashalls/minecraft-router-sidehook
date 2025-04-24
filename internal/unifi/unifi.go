package unifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type UnifiErrorResponse struct {
	Code      string                 `json:"code"`
	Details   map[string]interface{} `json:"details"`
	ErrorCode int                    `json:"errorCode"`
	Message   string                 `json:"message"`
}

type httpClient struct {
	Config *UnifiConfig
	*http.Client
	csrf       string
	ClientURLs *ClientURLs
}

type ClientURLs struct {
	Login            string
	NetworkObject    string
	PutNetworkObject string
}

const (
	unifiLoginPath                    = "%s/api/auth/login"
	unifiLoginPathExternal            = "%s/api/login"
	unifiNetworkObjectPath            = "%s/proxy/network/api/s/%s/rest/firewallgroup"
	unifiNetworkObjectPathExternal    = "%s/api/s/%s/rest/firewallgroup"
	unifiPutNetworkObjectPath         = "%s/proxy/network/api/s/%s/rest/firewallgroup/%s"
	unifiPutNetworkObjectPathExternal = "%s/proxy/network/api/s/%s/rest/firewallgroup/%s"
)

func NewClient(config *UnifiConfig) (*httpClient, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	client := &httpClient{
		Config: config,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SkipTLSVerify},
			},
			Jar: jar,
		},
		ClientURLs: &ClientURLs{
			Login:            unifiLoginPath,
			NetworkObject:    unifiNetworkObjectPath,
			PutNetworkObject: unifiPutNetworkObjectPath,
		},
	}

	if client.Config.ExternalController {
		log.Info("Heads up, you have external controller enabled. I have not tested this with non-UniFi controllers.")
		log.Info("If you have issues, please open an issue on GitHub.")
		client.ClientURLs.Login = unifiLoginPathExternal
		client.ClientURLs.NetworkObject = unifiNetworkObjectPathExternal
		client.ClientURLs.PutNetworkObject = unifiPutNetworkObjectPathExternal
	}

	if client.Config.ApiKey != "" {
		return client, nil
	}

	log.Info("UNIFI_USER and UNIFI_PASSWORD are deprecated, please switch to using UNIFI_API_KEY instead")

	if err := client.login(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *httpClient) login() error {
	jsonBody, err := json.Marshal(Login{
		Username: c.Config.User,
		Password: c.Config.Password,
		Remember: true,
	})
	if err != nil {
		return err
	}

	resp, err := c.doRequest(
		http.MethodPost,
		FormatUrl(c.ClientURLs.Login, c.Config.Host),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		log.Error("login failed", zap.String("status", resp.Status), zap.String("response", string(respBody)))
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	if csrf := resp.Header.Get("x-csrf-token"); csrf != "" {
		c.csrf = resp.Header.Get("x-csrf-token")
	}
	return nil
}

func (c *httpClient) doRequest(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Config.ApiKey == "" {
		if csrf := resp.Header.Get("X-CSRF-Token"); csrf != "" {
			c.csrf = csrf
		}
		// If the status code is 401, re-login and retry the request
		if resp.StatusCode == http.StatusUnauthorized {
			log.Debug("received 401 unauthorized, attempting to re-login")
			if err := c.login(); err != nil {
				log.Error("re-login failed", zap.Error(err))
				return nil, err
			}
			c.setHeaders(req)

			log.Debug("retrying request after re-login")

			resp, err = c.Client.Do(req)
			if err != nil {
				log.Error("Retry request failed", zap.Error(err))
				return nil, err
			}
		}
	}

	// It is unknown at this time if the UniFi API returns anything other than 200 for these types of requests.
	if resp.StatusCode != http.StatusOK {
		body, bodyErr := io.ReadAll(io.LimitReader(resp.Body, 512))
		if bodyErr != nil {
			return nil, bodyErr
		}

		var apiError UnifiErrorResponse
		if err := json.Unmarshal(body, &apiError); err != nil {
			return nil, fmt.Errorf("failed to decode json: %w", err)
		}

		return nil, fmt.Errorf("%s request to %s returned %d: %s", method, path, resp.StatusCode, apiError.Message)
	}

	return resp, nil
}

func (c *httpClient) setHeaders(req *http.Request) {
	if c.Config.ApiKey != "" {
		req.Header.Set("X-API-KEY", c.Config.ApiKey)
	} else {
		req.Header.Set("X-CSRF-Token", c.csrf)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
}

var (
	IPv4GroupType = "address-group"
	IPv6GroupType = "ipv6-address-group"
)

func (c *httpClient) CreateNetworkObject(object NetworkGroup) error {
	if object.GroupType != IPv4GroupType && object.GroupType != IPv6GroupType {
		return fmt.Errorf("invalid group type: %s", object.GroupType)
	}
	jsonStringify, err := json.Marshal(object)
	if err != nil {
		log.Error("failed to marshal group members", zap.Error(err))
		return err
	}

	resp, err := c.doRequest(
		http.MethodPost,
		FormatUrl(c.ClientURLs.NetworkObject, c.Config.Host, c.Config.Site),
		bytes.NewBuffer(jsonStringify),
	)
	if err != nil {
		log.Error("failed to create network object", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("failed to create network object", zap.String("status", resp.Status), zap.String("response", string(body)))
		return fmt.Errorf("failed to create network object: %s", resp.Status)
	}
	log.Info("created network object", zap.Any("object", object))
	return nil
}

// https://172.16.0.1/proxy/network/api/s/default/rest/firewallgroup
func (c *httpClient) FetchNetworkObjects() ([]NetworkGroup, error) {
	resp, err := c.doRequest(
		http.MethodGet,
		FormatUrl(c.ClientURLs.NetworkObject, c.Config.Host, c.Config.Site),
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var networkObject UnifiNetworkObject
	if err := json.NewDecoder(resp.Body).Decode(&networkObject); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return networkObject.Data, nil
}

// https://172.16.0.1/proxy/network/api/s/default/rest/firewallgroup/68096de988d7534d90ca2d52
func (c *httpClient) UpdateNetworkObject(object *NetworkGroup) error {
	jsonStringify, err := json.Marshal(object)
	if err != nil {
		log.Error("failed to marshal group members", zap.Error(err))
		return err
	}

	resp, err := c.doRequest(
		http.MethodPut,
		FormatUrl(c.ClientURLs.NetworkObject, c.Config.Host, c.Config.Site, object.ID),
		bytes.NewBuffer(jsonStringify),
	)
	if err != nil {
		log.Error("failed to update network object", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("failed to update network object", zap.String("status", resp.Status), zap.String("response", string(body)))
		return fmt.Errorf("failed to update network object: %s", resp.Status)
	}
	log.Info("updated network object", zap.Any("object", object))
	return nil
}
