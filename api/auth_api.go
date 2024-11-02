package api

import (
	"encoding/json"
	"guardian/internal/models"
	"guardian/internal/services"
	"guardian/utlis/logger"
	"net/http"
)

type AuthController struct {
	userService services.UserServiceInterface
}

func NewAuthController(userService services.UserServiceInterface) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

func (h *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.GetLogger().Errorf("Error:%v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.userService.Login(req)
	if err != nil {
		logger.GetLogger().Errorf("Error:%v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		logger.GetLogger().Errorf("Error:%v", err)
		http.Error(w, "Error encoding token", http.StatusInternalServerError)
		return
	}
}

func (h *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.GetLogger().Errorf("Error:%v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.userService.SignUp(req)
	if err != nil {
		logger.GetLogger().Errorf("Error:%v", err)
		http.Error(w, "Error in sign up", http.StatusInternalServerError)
		return
	}
}

func (h *AuthController) DeleteUser(_ http.ResponseWriter, _ *http.Request) {
}

func (h *AuthController) ActivateUser(_ http.ResponseWriter, _ *http.Request) {
}

func (h *AuthController) UpdateUser(_ http.ResponseWriter, _ *http.Request) {
}
