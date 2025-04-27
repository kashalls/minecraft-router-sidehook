package discord

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itzg/mc-router/server"
	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"go.uber.org/zap"
)

func InitServer(config DiscordConfig) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		var payload server.WebhookNotifierPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		log.Debug("Received webhook payload", zap.Any("payload", payload))

		params, err := BuildMessage(config.DiscordTemplate, payload)
		if err != nil {
			log.Error("Failed to build webhook message", zap.Error(err))
			http.Error(w, "Failed to build webhook message", http.StatusInternalServerError)
			return
		}

		log.Debug("Sending webhook message", zap.Any("params", params))

		webhookURL := config.WebhookURL
		if config.WebhookURLIsTemplate {
			var err error
			webhookURL, err = TemplateWebhookUrl(config.WebhookURL, payload)
			if err != nil {
				log.Error("Failed to template webhook URL", zap.Error(err))
				http.Error(w, "Failed to template webhook URL", http.StatusInternalServerError)
				return
			}
		}

		if err := SendWebhookMessage(webhookURL, params); err != nil {
			log.Error("Failed to send webhook message", zap.Error(err))
			http.Error(w, "Failed to send webhook message", http.StatusInternalServerError)
			return
		}

		log.Debug("Received audit webhook: %v", zap.Any("payload", payload))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Audit Webhook Received"))
	})
	return router
}
