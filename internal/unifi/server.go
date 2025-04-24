package unifi

import (
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"go.uber.org/zap"
)

func InitServer(unifi *httpClient) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		action := r.URL.Query().Get("action")

		if ip == "" || action == "" || (action != "add" && action != "remove") {
			http.Error(w, "Missing ip or action query param", http.StatusBadRequest)
			return
		}

		var objectType = unifi.Config.IPv4ObjectName
		parsedIp := net.ParseIP(ip)
		if parsedIp == nil {
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}

		if parsedIp.To16() != nil {
			objectType = unifi.Config.IPv6ObjectName
		}

		log.Info("recieved request for action", zap.String("action", action), zap.String("ip", ip), zap.String("objectType", objectType))
		// Log the action and IP address

		groups, err := unifi.FetchNetworkObjects()
		if err != nil {
			log.Error("error fetching network objects", zap.Error(err))
			http.Error(w, "Error fetching network objects", http.StatusInternalServerError)
			return
		}

		object := Find(groups, func(g NetworkGroup) bool {
			return g.Name == objectType
		})

		ipExists := Find(object.GroupMembers, func(objectIp string) bool {
			return ip == objectIp
		})

		if action == "add" && ipExists != nil {
			log.Info("tried to add ip that already exists in the group", zap.String("ip", ip))
			http.Error(w, "IP already exists in the group", http.StatusBadRequest)
			return
		}

		if action == "remove" && ipExists == nil {
			log.Info("tried to remove ip that does not exist in the group", zap.String("ip", ip))
			http.Error(w, "IP does not exist in the group", http.StatusBadRequest)
			return
		}

		if action == "add" {
			object.GroupMembers = append(object.GroupMembers, ip)
		} else if action == "remove" {
			for i, s := range object.GroupMembers {
				if s == ip {
					object.GroupMembers = append(object.GroupMembers[:i], object.GroupMembers[i+1:]...)
					break
				}
			}
		}

		if err := unifi.UpdateNetworkObject(object); err != nil {
			log.Error("error updating network object", zap.Error(err))
			http.Error(w, "Error updating network object", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Network object updated successfully"))

	})
	return router
}
