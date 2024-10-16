package models

import (
	"github.com/google/uuid"
)

type SendRequest struct {
	UserID uuid.UUID     `json:"user_id"`
	ChatID *uuid.UUID    `json:"chat_id,omitempty"`
	Prompt string        `json:"prompt"`
	Target []TargetModel `json:"targets"`
}

type SendResponse struct {
	Status string  `json:"status"`
	Target AIModel `json:"target,omitempty"`
}

type AIModel struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type RefereeModel struct {
	ModelID uuid.UUID `json:"model_id"`
	Token   string    `json:"token,omitempty"`
}

type TargetModel struct {
	ModelID uuid.UUID `json:"model_id"`
	Token   string    `json:"token"`
}

type Usage struct {
	UserID                 uuid.UUID `json:"user_id"`
	TargetModelID          uuid.UUID `json:"target_model_id"`
	InputTokenConsumption  int       `json:"input_token_consumption"`
	OutputTokenConsumption int       `json:"output_token_consumption"`
}
