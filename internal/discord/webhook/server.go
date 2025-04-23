package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kashalls/minecraft-router-sidehook/cmd/discord/configuration"
	"github.com/kashalls/minecraft-router-sidehook/internal/constants"
	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"go.uber.org/zap"
)

func InitServer() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/webhook", HandleAudit)

	return router

}

func HandleAudit(w http.ResponseWriter, r *http.Request) {
	var payload constants.WebhookNotifierPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Debug("Received webhook payload", zap.Any("payload", payload))

	params, err := BuildMessage(payload)
	if err != nil {
		log.Error("Failed to build webhook message", zap.Error(err))
		http.Error(w, "Failed to build webhook message", http.StatusInternalServerError)
		return
	}

	log.Debug("Sending webhook message", zap.Any("params", params))

	if err := SendWebhookMessage(configuration.Config.Webhook, params); err != nil {
		log.Error("Failed to send webhook message", zap.Error(err))
		http.Error(w, "Failed to send webhook message", http.StatusInternalServerError)
		return
	}

	log.Debug("Received audit webhook: %v", zap.Any("payload", payload))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Audit Webhook Received"))
}
