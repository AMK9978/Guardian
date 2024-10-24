package api

import (
	"guardian/internal/models"
	"guardian/internal/services"

	"encoding/json"
	"net/http"
)

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

