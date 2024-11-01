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

// RefereeRequest represents a request sending to the referee plugins
type RefereeRequest struct {
	UserID   primitive.ObjectID `json:"user_id"`
	Chat     string             `json:"chat,omitempty"`
	Prompt   string             `json:"prompt"`
	TargetID primitive.ObjectID `json:"target_id"`
}

// SendResponse represents the response from a send operation.
type SendResponse struct {
	Status bool `json:"status"`
}
