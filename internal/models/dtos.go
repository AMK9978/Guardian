package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"guardian/internal/models/entities"
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
	UserID primitive.ObjectID   `json:"user_id"`
	ChatID *primitive.ObjectID  `json:"chat_id,omitempty"`
	Prompt string               `json:"prompt"`
	Target entities.TargetModel `json:"target"`
}

// SendResponse represents the response from a send operation.
type SendResponse struct {
	Success bool                 `json:"success"`
	Status  string               `json:"status"`
	Target  entities.TargetModel `json:"target,omitempty"`
}
