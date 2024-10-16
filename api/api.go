package api

import (
	"encoding/json"
	"guardian/internal/models"
	"net/http"

	"guardian/internal/services"
	"guardian/utlis/logger"
)

func SendHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SendRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	status, err := services.ProcessPrompt(req)
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
