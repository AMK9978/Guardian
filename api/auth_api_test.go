package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"guardian/internal/mocks"
	"guardian/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

)

var (
	ErrLogin  = errors.New("login error")
	ErrSignup = errors.New("signup error")
)

func TestAuthController_Login(t *testing.T) {
	t.Parallel()

	mockService := new(mocks.MockUserService)
	controller := NewAuthController(mockService)

	reqBody := models.LoginRequest{
		Email:    "test@test.com",
		Password: "test",
	}
	token := "sample"

	t.Run("successful login", func(t *testing.T) {
		t.Parallel()

		mockService.On("Login", reqBody).Return(token, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/login", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.Login(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var respBody map[string]string
		_ = json.NewDecoder(rec.Body).Decode(&respBody)
		assert.Equal(t, token, respBody["token"])
		mockService.AssertCalled(t, "Login", reqBody)
		mockService.On("Login", reqBody).Unset()
	})

	t.Run("login with error", func(t *testing.T) {
		t.Parallel()
		mockService := new(mocks.MockUserService)
		controller := NewAuthController(mockService)
		mockService.On("Login", reqBody).Return("", ErrLogin)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/login", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.Login(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertCalled(t, "Login", reqBody)
	})
}

func TestAuthController_SignUp(t *testing.T) {
	t.Parallel()
	mockService := new(mocks.MockUserService)
	controller := NewAuthController(mockService)

	reqBody := models.SignUpRequest{
		Name:     "test",
		Email:    "test@test.com",
		Password: "test",
	}

	t.Run("Signup successfully", func(t *testing.T) {
		t.Parallel()
		mockService.On("SignUp", reqBody).Return(nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.SignUp(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertCalled(t, "SignUp", reqBody)
	})

	mockService.ExpectedCalls = nil
	t.Run("Signup fails", func(t *testing.T) {
		t.Parallel()
		mockService := new(mocks.MockUserService)
		controller := NewAuthController(mockService)
		mockService.On("SignUp", reqBody).Return(ErrSignup)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.SignUp(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertCalled(t, "SignUp", reqBody)
	})
}
