package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndreyKuskov2/gophermart/internal/config"
	"github.com/AndreyKuskov2/gophermart/internal/models"
	"github.com/AndreyKuskov2/gophermart/internal/storage"
	"github.com/AndreyKuskov2/gophermart/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGophermartUserServicer is a mock implementation of GophermartUserServicer
type MockGophermartUserServicer struct {
	mock.Mock
}

func (m *MockGophermartUserServicer) RegisterUserService(ctx context.Context, user models.UserCreditials) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockGophermartUserServicer) GetUserService(ctx context.Context, user models.UserCreditials) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func getTestConfig() *config.Config {
	return &config.Config{
		JWTSecretToken: "test-secret",
	}
}

func getTestLogger() *logger.Logger {
	log, _ := logger.NewLogger()
	return log
}

func TestRegisterUserHandler_Success(t *testing.T) {
	mockService := &MockGophermartUserServicer{}
	cfg := getTestConfig()
	log := getTestLogger()
	h := NewGophermartUserHandlers(mockService, cfg, log)

	user := models.UserCreditials{Login: "testuser", Password: "testpass"}
	mockService.On("RegisterUserService", mock.Anything, user).Return(1, nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RegisterUserHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("Authorization"))
	mockService.AssertExpectations(t)
}

func TestRegisterUserHandler_BadRequest(t *testing.T) {
	mockService := &MockGophermartUserServicer{}
	cfg := getTestConfig()
	log := getTestLogger()
	h := NewGophermartUserHandlers(mockService, cfg, log)

	// Missing password
	user := models.UserCreditials{Login: "testuser"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RegisterUserHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRegisterUserHandler_Conflict(t *testing.T) {
	mockService := &MockGophermartUserServicer{}
	cfg := getTestConfig()
	log := getTestLogger()
	h := NewGophermartUserHandlers(mockService, cfg, log)

	user := models.UserCreditials{Login: "testuser", Password: "testpass"}
	mockService.On("RegisterUserService", mock.Anything, user).Return(0, storage.ErrUserIsExist)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RegisterUserHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestRegisterUserHandler_InternalServerError_Service(t *testing.T) {
	mockService := &MockGophermartUserServicer{}
	cfg := getTestConfig()
	log := getTestLogger()
	h := NewGophermartUserHandlers(mockService, cfg, log)

	user := models.UserCreditials{Login: "testuser", Password: "testpass"}
	mockService.On("RegisterUserService", mock.Anything, user).Return(0, errors.New("db error"))

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RegisterUserHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestLoginUserHandler_Success(t *testing.T) {
	mockService := &MockGophermartUserServicer{}
	cfg := getTestConfig()
	log := getTestLogger()
	h := NewGophermartUserHandlers(mockService, cfg, log)

	user := models.UserCreditials{Login: "testuser", Password: "testpass"}
	mockService.On("GetUserService", mock.Anything, user).Return(1, nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.LoginUserHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("Authorization"))
	mockService.AssertExpectations(t)
}

func TestLoginUserHandler_BadRequest(t *testing.T) {
	mockService := &MockGophermartUserServicer{}
	cfg := getTestConfig()
	log := getTestLogger()
	h := NewGophermartUserHandlers(mockService, cfg, log)

	user := models.UserCreditials{Login: "testuser"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.LoginUserHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLoginUserHandler_Unauthorized(t *testing.T) {
	mockService := &MockGophermartUserServicer{}
	cfg := getTestConfig()
	log := getTestLogger()
	h := NewGophermartUserHandlers(mockService, cfg, log)

	user := models.UserCreditials{Login: "testuser", Password: "wrongpass"}
	mockService.On("GetUserService", mock.Anything, user).Return(0, sql.ErrNoRows)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.LoginUserHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestLoginUserHandler_Unauthorized_InvalidData(t *testing.T) {
	mockService := &MockGophermartUserServicer{}
	cfg := getTestConfig()
	log := getTestLogger()
	h := NewGophermartUserHandlers(mockService, cfg, log)

	user := models.UserCreditials{Login: "testuser", Password: "wrongpass"}
	mockService.On("GetUserService", mock.Anything, user).Return(0, storage.ErrInvalidData)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.LoginUserHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestLoginUserHandler_InternalServerError_Service(t *testing.T) {
	mockService := &MockGophermartUserServicer{}
	cfg := getTestConfig()
	log := getTestLogger()
	h := NewGophermartUserHandlers(mockService, cfg, log)

	user := models.UserCreditials{Login: "testuser", Password: "testpass"}
	mockService.On("GetUserService", mock.Anything, user).Return(0, errors.New("db error"))

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.LoginUserHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}
