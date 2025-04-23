package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/kashalls/minecraft-router-sidehook/cmd/discord/configuration"
	"github.com/kashalls/minecraft-router-sidehook/internal/constants"
)

var DefaultTemplate = `{
	"username": "{{.Event}}",
	"content": "New Event: {{.Status}}"
}`

func BuildMessage(data constants.WebhookNotifierPayload) (*discordgo.WebhookParams, error) {
	tmplStr := configuration.Config.WebhookTemplate
	if tmplStr == "" {
		tmplStr = DefaultTemplate
	}

	tmpl, err := template.New("webhook").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	var params discordgo.WebhookParams
	if err := json.Unmarshal(buf.Bytes(), &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook params: %w", err)
	}

	return &params, nil
}

func SendWebhookMessage(url string, message *discordgo.WebhookParams) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "MinecraftRouterSidehook/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Discord webhook returned status: %s", resp.Status)
	}
	return nil
}
