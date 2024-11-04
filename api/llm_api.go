package api

import (
	"encoding/json"
	"fmt"
	"guardian/internal/middleware"
	"guardian/internal/models"
	"guardian/internal/services"
	"guardian/utlis/logger"
	"io"
	"net/http"
)

type SendHandlerController struct {
	promptService      services.PromptServiceInterface
	targetModelService services.TargetModelServiceInterface
	middleware         middleware.Interface
}

func NewSendHandlerController(promptService services.PromptServiceInterface,
	targetModelService services.TargetModelServiceInterface, m middleware.Interface,
) *SendHandlerController {
	return &SendHandlerController{
		promptService:      promptService,
		targetModelService: targetModelService,
		middleware:         m,
	}
}

func (h *SendHandlerController) SendHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.middleware.GetUserFromContext(r)
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}

	var reqBody models.PluginRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		logger.GetLogger().Errorf("error in sendhandler %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	reqBody.UserID = *userID
	targetLLM, err := h.targetModelService.GetTargetModel(r.Context(), reqBody.TargetID)
	if err != nil {
		logger.GetLogger().Errorf("error in resolving the target LLM %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	result, err := h.promptService.ProcessPrompt(r.Context(), &reqBody)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !result {
		resp := models.PluginResponse{Status: result}
		w.Header().Set("Content-Type", "application/json")
		respBody, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(respBody)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	newReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetLLM.Address, r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	for k, v := range r.Header {
		newReq.Header[k] = v
	}

	resp, err := h.promptService.SendPrompt(r.Context(), newReq)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.returnResponseToUser(w, resp)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *SendHandlerController) returnResponseToUser(w http.ResponseWriter, resp *http.Response) error {
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

// TODO: selectTargetLLM should choose if the user/group let the system choose the appropriate Target LLM.
// func selectTargetLLM() {
//}
