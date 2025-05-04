package handlers

import (
	"base_structure/src/api/dto"
	"base_structure/src/api/validations"
	"base_structure/src/config"
	"base_structure/src/pkg/service_errors"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// -----------------------------------------------------------------------------
// Test bootstrap

func TestMain(m *testing.M) {
	// register custom validators (same as server boot)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("ir_mobile", validations.IranianMobileNumberValidator, true)
		_ = v.RegisterValidation("password", validations.PasswordValidator, true)
	}
	_ = os.Setenv("APP_ENV", "development")
	os.Exit(m.Run())
}

// -----------------------------------------------------------------------------
// Mock implementing handlers.UserSvcPort

type mockUserService struct{ mock.Mock }

func (m *mockUserService) SendOtp(r *dto.GetOtpRequest) error {
	return m.Called(r).Error(0)
}
func (m *mockUserService) LoginByUsername(r *dto.LoginByUsernameRequest) (*dto.TokenDetail, error) {
	args := m.Called(r)
	td, _ := args.Get(0).(*dto.TokenDetail)
	return td, args.Error(1)
}
func (m *mockUserService) RegisterByUsername(r *dto.RegisterByUsernameRequest) error {
	return m.Called(r).Error(0)
}
func (m *mockUserService) RegisterLoginByMobileNumber(r *dto.RegisterLoginByMobileRequest) (*dto.TokenDetail, error) {
	args := m.Called(r)
	td, _ := args.Get(0).(*dto.TokenDetail)
	return td, args.Error(1)
}

// -----------------------------------------------------------------------------
// Helpers for SendOtp

func newSendOtpEnv() (*gin.Engine, *mockUserService) {
	ms := new(mockUserService)
	h := NewUserHandlerWithSvc(&config.Config{}, ms)
	r := gin.New()
	r.POST("/api/v1/users/send-otp", h.SendOtp)
	return r, ms
}

func perform(r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// -----------------------------------------------------------------------------
// SendOtp tests

func TestSendOtp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("happy path", func(t *testing.T) {
		r, ms := newSendOtpEnv()
		inp := dto.GetOtpRequest{MobileNumber: "09123456789"}
		ms.On("SendOtp", mock.AnythingOfType("*dto.GetOtpRequest")).Return(nil).Once()

		w := perform(r, "/api/v1/users/send-otp", inp)
		assert.Equal(t, http.StatusCreated, w.Code)
		ms.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		r, _ := newSendOtpEnv()
		w := perform(r, "/api/v1/users/send-otp", map[string]string{"mobileNumber": "123"})
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("otp exists", func(t *testing.T) {
		r, ms := newSendOtpEnv()
		ms.On("SendOtp", mock.AnythingOfType("*dto.GetOtpRequest")).
			Return(&service_errors.ServiceError{EndUserMessage: service_errors.OtpExists}).Once()

		w := perform(r, "/api/v1/users/send-otp", dto.GetOtpRequest{MobileNumber: "09123456789"})
		assert.Equal(t, http.StatusConflict, w.Code)
		ms.AssertExpectations(t)
	})
}

// -----------------------------------------------------------------------------
// Helpers for LoginByUsername

func newLoginEnv() (*gin.Engine, *mockUserService) {
	ms := new(mockUserService)
	h := NewUserHandlerWithSvc(&config.Config{}, ms)
	r := gin.New()
	r.POST("/api/v1/users/login-by-username", h.LoginByUsername)
	return r, ms
}

// -----------------------------------------------------------------------------
// LoginByUsername tests

func TestLoginByUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("happy path", func(t *testing.T) {
		r, ms := newLoginEnv()
		inp := dto.LoginByUsernameRequest{Username: "admin", Password: "secret"}
		td := &dto.TokenDetail{AccessToken: "at", RefreshToken: "rt"}
		ms.On("LoginByUsername", mock.AnythingOfType("*dto.LoginByUsernameRequest")).Return(td, nil).Once()

		w := perform(r, "/api/v1/users/login-by-username", inp)
		assert.Equal(t, http.StatusCreated, w.Code)
		ms.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		r, _ := newLoginEnv()
		// username too short -> should fail binding
		w := perform(r, "/api/v1/users/login-by-username", dto.LoginByUsernameRequest{Username: "ad", Password: "123"})
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		r, ms := newLoginEnv()
		ms.On("LoginByUsername", mock.Anything).
			Return((*dto.TokenDetail)(nil),
				&service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}).Once()
		w := perform(r, "/api/v1/users/login-by-username", dto.LoginByUsernameRequest{Username: "admin", Password: "badpass"})
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		ms.AssertExpectations(t)
	})
}
