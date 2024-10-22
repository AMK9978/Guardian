package api

import (
	"encoding/json"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"guardian/internal/models"
	"guardian/internal/services"
	"guardian/utlis/logger"
	"net/http"
)

type SendHandlerController struct {
	promptService services.PromptServiceInterface
}

func NewSendHandlerController(promptService *services.PromptService) *SendHandlerController {
	return &SendHandlerController{
		promptService: promptService,
	}
}

type AuthController struct {
	userService services.UserServiceInterface
}

func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

func (h *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.userService.Login(req)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthController) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthController) ActivateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthController) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *SendHandlerController) SendHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	userIDStr := claims["user_id"].(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		logger.GetLogger().Fatalf("error in reading the userID from the user's token: %s", userIDStr)
	}

	var req models.SendRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.GetLogger().Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.UserID = userID

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
