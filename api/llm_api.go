package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"guardian/internal/middleware"
	"guardian/internal/models"
	"guardian/internal/services"
	"guardian/utlis/logger"
)

type SendHandlerController struct {
	promptService      services.PromptServiceInterface
	targetModelService *services.TargetModelService
}

func NewSendHandlerController(promptService *services.PromptService,
	targetModelService *services.TargetModelService,
) *SendHandlerController {
	return &SendHandlerController{
		promptService:      promptService,
		targetModelService: targetModelService,
	}
}

func (h *SendHandlerController) SendHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserFromContext(r)
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}

	var reqBody models.RefereeRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		logger.GetLogger().Errorf("error in sendhandler %v", err)
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	reqBody.UserID = *userID
	targetLLM, err := h.targetModelService.GetTargetModel(reqBody.TargetID)
	if err != nil {
		logger.GetLogger().Errorf("error in resolving the target LLM %v", err)
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	result, err := h.promptService.ProcessPrompt(r.Context(), &reqBody)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !result {
		resp := models.SendResponse{Status: result}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	newReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetLLM.Address, r.Body)
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
