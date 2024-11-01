package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"guardian/internal/mocks"
	"guardian/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestAuthController_Login(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := NewAuthController(mockService)

	reqBody := models.LoginRequest{
		Email:    "test@test.com",
		Password: "test",
	}
	token := "sample"

	t.Run("successful login", func(t *testing.T) {
		mockService.On("Login", reqBody).Return(token, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.Login(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var respBody map[string]string
		json.NewDecoder(rec.Body).Decode(&respBody)
		assert.Equal(t, token, respBody["token"])
		mockService.AssertCalled(t, "Login", reqBody)
	})

	mockService.ExpectedCalls = nil
	t.Run("login with error", func(t *testing.T) {
		mockService.On("Login", reqBody).Return("", errors.New("login error"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.Login(rec, req)
		fmt.Println(rec.Code)
		fmt.Println(rec)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertCalled(t, "Login", reqBody)
	})
}

func TestAuthController_SignUp(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := NewAuthController(mockService)

	reqBody := models.SignUpRequest{
		Name:     "test",
		Email:    "test@test.com",
		Password: "test",
	}

	t.Run("Signup successfully", func(t *testing.T) {
		mockService.On("SignUp", reqBody).Return(nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.SignUp(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertCalled(t, "SignUp", reqBody)
	})

	mockService.ExpectedCalls = nil
	t.Run("Signup fails", func(t *testing.T) {
		mockService.On("SignUp", reqBody).Return(errors.New("signup error"))

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.SignUp(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertCalled(t, "SignUp", reqBody)
	})
}
