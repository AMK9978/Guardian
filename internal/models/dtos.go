package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SendRequest represents a request to send a prompt.
type SendRequest struct {
	ChatID   *primitive.ObjectID `json:"chat_id,omitempty"`
	Prompt   string              `json:"prompt"`
	TargetID primitive.ObjectID  `json:"target_id"`
}

// PluginRequest represents a request sending to the referee plugins
type PluginRequest struct {
	UserID   primitive.ObjectID `json:"user_id"`
	Chat     string             `json:"chat,omitempty"`
	Address  string             `json:"address,omitempty"`
	Prompt   string             `json:"prompt"`
	TargetID primitive.ObjectID `json:"target_id"`
}

// PluginResponse represents the response from a send operation.
type PluginResponse struct {
	Status bool   `json:"status"`
	Score  uint32 `json:"score,omitempty"`
}
