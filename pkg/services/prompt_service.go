package services

import (
	"errors"

	"guardian/models"
)

func ProcessPrompt(req models.SendRequest) (string, error) {
	if req.Prompt == "" {
		return "", errors.New("empty prompt")
	}

	if pipeline(req) {
		return "malicious", nil
	}

	return "benign", nil
}

func pipeline(req models.SendRequest) bool {

	return false
}
