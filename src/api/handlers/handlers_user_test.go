package handlers

import (
	"base_structure/src/api/dto"
	"base_structure/src/api/validations"
	"base_structure/src/config"
	"base_structure/src/constants"
	"base_structure/src/pkg/service_errors"
	"bytes"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

/* ------------------------------------------------------------------------- */
/* Test bootstrap                                                            */

func TestMain(m *testing.M) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("ir_mobile", validations.IranianMobileNumberValidator, true)
		_ = v.RegisterValidation("password", validations.PasswordValidator, true)
	}
	_ = os.Setenv("APP_ENV", "development")
	os.Exit(m.Run())
}

/* ------------------------------------------------------------------------- */
/* Mock service (implements services.UserServiceIface)                       */

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

/* ------------------------------------------------------------------------- */
/* Generic HTTP helpers                                                      */

func performWithHdr(r *gin.Engine, path string, body any, hdr map[string]string) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func perform(r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	return performWithHdr(r, path, body, nil)
}

/* ------------------------------------------------------------------------- */
/* Router builders per handler                                               */

func newSendOtpEnv() (*gin.Engine, *mockUserService) {
	ms := new(mockUserService)
	h := NewUserHandlerWithSvc(&config.Config{}, ms)
	r := gin.New()
	r.POST("/api/v1/users/send-otp", h.SendOtp)
	return r, ms
}
func newLoginEnv() (*gin.Engine, *mockUserService) {
	ms := new(mockUserService)
	h := NewUserHandlerWithSvc(&config.Config{}, ms)
	r := gin.New()
	r.POST("/api/v1/users/login-by-username", h.LoginByUsername)
	return r, ms
}
func newRegisterEnv() (*gin.Engine, *mockUserService) {
	ms := new(mockUserService)
	h := NewUserHandlerWithSvc(&config.Config{}, ms)
	r := gin.New()
	r.POST("/api/v1/users/register-by-username", h.RegisterByUsername)
	return r, ms
}
func newMobileEnv() (*gin.Engine, *mockUserService) {
	ms := new(mockUserService)
	h := NewUserHandlerWithSvc(&config.Config{}, ms)
	r := gin.New()
	r.POST("/api/v1/users/login-by-mobile", h.RegisterLoginByMobileNumber)
	return r, ms
}
func newLogoutEnv(cfg *config.Config) *gin.Engine {
	h := NewUserHandlerWithSvc(cfg, new(mockUserService))
	r := gin.New()
	r.POST("/api/v1/users/logout", h.Logout)
	return r
}

/* ------------------------------------------------------------------------- */
/* SendOtp                                                                   */

func TestSendOtp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("happy path", func(t *testing.T) {
		r, ms := newSendOtpEnv()
		ms.On("SendOtp", mock.Anything).Return(nil).Once()

		w := perform(r, "/api/v1/users/send-otp", dto.GetOtpRequest{MobileNumber: "09123456789"})
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
		ms.On("SendOtp", mock.Anything).
			Return(&service_errors.ServiceError{EndUserMessage: service_errors.OtpExists}).Once()

		w := perform(r, "/api/v1/users/send-otp", dto.GetOtpRequest{MobileNumber: "09123456789"})
		assert.Equal(t, http.StatusConflict, w.Code)
		ms.AssertExpectations(t)
	})
}

/* ------------------------------------------------------------------------- */
/* LoginByUsername                                                           */

func TestLoginByUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("happy path", func(t *testing.T) {
		r, ms := newLoginEnv()
		td := &dto.TokenDetail{AccessToken: "at", RefreshToken: "rt"}
		ms.On("LoginByUsername", mock.Anything).Return(td, nil).Once()

		w := perform(r, "/api/v1/users/login-by-username",
			dto.LoginByUsernameRequest{Username: "admin", Password: "secret"})
		assert.Equal(t, http.StatusCreated, w.Code)
		ms.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		r, _ := newLoginEnv()
		bad := dto.LoginByUsernameRequest{Username: "ad", Password: "123"}
		w := perform(r, "/api/v1/users/login-by-username", bad)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("invalid creds", func(t *testing.T) {
		r, ms := newLoginEnv()
		ms.On("LoginByUsername", mock.Anything).
			Return((*dto.TokenDetail)(nil),
				&service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}).Once()

		w := perform(r, "/api/v1/users/login-by-username",
			dto.LoginByUsernameRequest{Username: "admin", Password: "badpass"})
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		ms.AssertExpectations(t)
	})
}

/* ------------------------------------------------------------------------- */
/* RegisterByUsername                                                        */

func TestRegisterByUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	valid := dto.RegisterByUsernameRequest{
		FirstName: "Ali", LastName: "Admin", Username: "admin123",
		Email: "a@a.com", Password: "P@ssw0rd",
	}

	t.Run("happy path", func(t *testing.T) {
		r, ms := newRegisterEnv()
		ms.On("RegisterByUsername", mock.Anything).Return(nil).Once()

		w := perform(r, "/api/v1/users/register-by-username", valid)
		assert.Equal(t, http.StatusCreated, w.Code)
		ms.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		r, _ := newRegisterEnv()
		bad := dto.RegisterByUsernameRequest{FirstName: "A", LastName: "B"}
		w := perform(r, "/api/v1/users/register-by-username", bad)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("username exists", func(t *testing.T) {
		r, ms := newRegisterEnv()
		ms.On("RegisterByUsername", mock.Anything).
			Return(&service_errors.ServiceError{EndUserMessage: service_errors.UsernameExists}).Once()
		w := perform(r, "/api/v1/users/register-by-username", valid)
		assert.Equal(t, http.StatusConflict, w.Code)
		ms.AssertExpectations(t)
	})

	t.Run("email exists", func(t *testing.T) {
		r, ms := newRegisterEnv()
		ms.On("RegisterByUsername", mock.Anything).
			Return(&service_errors.ServiceError{EndUserMessage: service_errors.EmailExists}).Once()
		w := perform(r, "/api/v1/users/register-by-username", valid)
		assert.Equal(t, http.StatusConflict, w.Code)
		ms.AssertExpectations(t)
	})
}

/* ------------------------------------------------------------------------- */
/* RegisterLoginByMobileNumber                                               */

func TestRegisterLoginByMobileNumber(t *testing.T) {
	gin.SetMode(gin.TestMode)

	valid := dto.RegisterLoginByMobileRequest{MobileNumber: "09123456789", Otp: "123456"}
	td := &dto.TokenDetail{AccessToken: "at", RefreshToken: "rt"}

	t.Run("happy path", func(t *testing.T) {
		r, ms := newMobileEnv()
		ms.On("RegisterLoginByMobileNumber", mock.Anything).Return(td, nil).Once()
		w := perform(r, "/api/v1/users/login-by-mobile", valid)
		assert.Equal(t, http.StatusCreated, w.Code)
		ms.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		r, _ := newMobileEnv()
		bad := dto.RegisterLoginByMobileRequest{MobileNumber: "123", Otp: "1"}
		w := perform(r, "/api/v1/users/login-by-mobile", bad)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("invalid otp", func(t *testing.T) {
		r, ms := newMobileEnv()
		ms.On("RegisterLoginByMobileNumber", mock.Anything).
			Return((*dto.TokenDetail)(nil),
				&service_errors.ServiceError{EndUserMessage: service_errors.OtpNotValid}).Once()
		w := perform(r, "/api/v1/users/login-by-mobile", valid)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		ms.AssertExpectations(t)
	})
}

/* ------------------------------------------------------------------------- */
/* Logout – uses real token/redis with miniredis                             */

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// start in‑memory Redis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("cannot start miniredis: %v", err)
	}
	defer mr.Close()

	host, port, _ := strings.Cut(mr.Addr(), ":")
	cfg := &config.Config{
		Jwt: config.JwtConfig{
			Secret: "secret", RefreshSecret: "refresh",
			AccessTokenExpireDuration: 60, RefreshTokenExpireDuration: 60,
		},
		Redis: config.RedisConfig{Host: host, Port: port},
	}

	// helper: create valid access/refresh pair
	makeToken := func() (string, string) {
		atClaims := jwt.MapClaims{
			constants.UserIdKey:     1,
			constants.ExpireTimeKey: time.Now().Add(time.Hour).Unix(),
		}
		rtClaims := jwt.MapClaims{
			constants.UserIdKey:     1,
			constants.ExpireTimeKey: time.Now().Add(time.Hour).Unix(),
		}
		at, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString([]byte(cfg.Jwt.Secret))
		rt, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims).SignedString([]byte(cfg.Jwt.RefreshSecret))
		return at, rt
	}

	t.Run("validation error (missing refresh)", func(t *testing.T) {
		r := newLogoutEnv(cfg)
		at, _ := makeToken()
		h := map[string]string{constants.AuthorizationHeaderKey: "Bearer " + at}
		w := performWithHdr(r, "/api/v1/users/logout",
			map[string]string{"refreshToken": ""}, h)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("missing auth header", func(t *testing.T) {
		r := newLogoutEnv(cfg)

		// create a valid refresh token but DO NOT set Authorization header
		_, rt := makeToken()

		w := perform(r,
			"/api/v1/users/logout",
			map[string]string{"refreshToken": rt})

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
