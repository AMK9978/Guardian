package api

import (
	"encoding/json"
	"fmt"
	"guardian/internal/middleware"
	"io"
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

	result, err := h.promptService.ProcessPrompt(req, r)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !result {
		resp := models.SendResponse{Success: result}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	newReq, err := http.NewRequestWithContext(r.Context(), r.Method, req.Target.Address, r.Body)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	for k, v := range r.Header {
		newReq.Header[k] = v
	}

	client := &http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = forwardResponseToUser(w, resp)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func forwardResponseToUser(w http.ResponseWriter, resp *http.Response) error {
	w.WriteHeader(resp.StatusCode)

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	_, err := io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write response body: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func selectTargetLLM() {

}
