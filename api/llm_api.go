package api

import (
	"encoding/json"
	"guardian/internal/middleware"
	"net/http"

	"guardian/internal/models"
	"guardian/internal/services"
	"guardian/utlis/logger"
)

type SendHandlerController struct {
	promptService services.PromptServiceInterface
}

func NewSendHandlerController(promptService *services.PromptService) *SendHandlerController {
	return &SendHandlerController{
		promptService: promptService,
	}
}

func (h *SendHandlerController) SendHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserFromContext(r)
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}

	var req models.SendRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.UserID = *userID

	status, err := h.promptService.ProcessPrompt(req)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := models.SendResponse{Status: status}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
