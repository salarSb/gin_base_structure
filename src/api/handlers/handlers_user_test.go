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

func TestMain(m *testing.M) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("ir_mobile", validations.IranianMobileNumberValidator, true)
		_ = v.RegisterValidation("password", validations.PasswordValidator, true)
	}
	_ = os.Setenv("APP_ENV", "development")
	os.Exit(m.Run())
}

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

func newTestEnv() (*gin.Engine, *mockUserService) {
	ms := new(mockUserService)
	h := NewUserHandlerWithSvc(&config.Config{}, ms)
	r := gin.New()
	r.POST("/api/v1/users/send-otp", h.SendOtp)
	return r, ms
}

func perform(r *gin.Engine, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/send-otp", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSendOtp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("happy path", func(t *testing.T) {
		r, ms := newTestEnv()
		inp := dto.GetOtpRequest{MobileNumber: "09123456789"}
		ms.On("SendOtp", mock.AnythingOfType("*dto.GetOtpRequest")).Return(nil).Once()
		w := perform(r, inp)
		assert.Equal(t, http.StatusCreated, w.Code)
		ms.AssertExpectations(t)
	})
	t.Run("validation error", func(t *testing.T) {
		r, _ := newTestEnv()
		w := perform(r, map[string]string{"mobileNumber": "123"})
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("otp exists", func(t *testing.T) {
		r, ms := newTestEnv()
		ms.On("SendOtp", mock.AnythingOfType("*dto.GetOtpRequest")).
			Return(&service_errors.ServiceError{EndUserMessage: service_errors.OtpExists}).Once()
		w := perform(r, dto.GetOtpRequest{MobileNumber: "09123456789"})
		assert.Equal(t, http.StatusConflict, w.Code)
		ms.AssertExpectations(t)
	})
}
