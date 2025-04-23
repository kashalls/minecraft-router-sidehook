package constants

import (
	"time"

	"github.com/google/uuid"
)

type WebhookNotifierPayload struct {
	Event           string      `json:"event"`
	Timestamp       time.Time   `json:"timestamp"`
	Status          string      `json:"status"`
	Client          *ClientInfo `json:"client"`
	Server          string      `json:"server"`
	PlayerInfo      *PlayerInfo `json:"player,omitempty"`
	BackendHostPort string      `json:"backend,omitempty"`
	Error           string      `json:"error,omitempty"`
}

type ClientInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type PlayerInfo struct {
	Name string    `json:"name"`
	Uuid uuid.UUID `json:"uuid"`
}
