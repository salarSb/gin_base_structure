package handlers

import (
	"base_structure/src/api/dto"
	"base_structure/src/api/helper"
	"base_structure/src/config"
	"base_structure/src/constants"
	"base_structure/src/pkg/service_errors"
	"base_structure/src/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

type UserHandler struct {
	cfg         *config.Config
	userService services.UserServiceIface
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{
		cfg:         cfg,
		userService: services.NewUserService(cfg),
	}
}

func NewUserHandlerWithSvc(cfg *config.Config, svc services.UserServiceIface) *UserHandler {
	return &UserHandler{
		cfg:         cfg,
		userService: svc,
	}
}

// SendOtp
// @Summary      Send OTP
// @Description  Sends a 6-digit one-time password (OTP) to a verified Iranian mobile number.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.GetOtpRequest  true  "Mobile number payload"
// @Success      201      {object}  helper.BaseHttpResponse                "OTP sent"
// @Failure      422      {object}  helper.BaseHttpResponse                "Validation error"
// @Failure      429      {object}  helper.BaseHttpResponse                "Rate-limited – too many requests"
// @Failure      500      {object}  helper.BaseHttpResponse                "Internal server error"
// @Router       /api/v1/users/send-otp [post]
func (h *UserHandler) SendOtp(c *gin.Context) {
	req := new(dto.GetOtpRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
		return
	}
	err = h.userService.SendOtp(req)
	if err != nil {
		c.AbortWithStatusJSON(
			helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
		return
	}
	// Call internal sms service
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse("otp sent", true, helper.Success))
}

// LoginByUsername
// @Summary      Login with username/password
// @Description  Authenticates a user by username & password and returns a JWT access/refresh pair.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LoginByUsernameRequest  true  "Credentials"
// @Success      201      {object}  helper.BaseHttpResponse{result=dto.TokenDetail}  "Tokens returned"
// @Failure      401      {object}  helper.BaseHttpResponse                            "Invalid credentials"
// @Failure      422      {object}  helper.BaseHttpResponse                            "Validation error"
// @Failure      500      {object}  helper.BaseHttpResponse                            "Internal server error"
// @Router       /api/v1/users/login-by-username [post]
func (h *UserHandler) LoginByUsername(c *gin.Context) {
	req := new(dto.LoginByUsernameRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
		return
	}
	token, err := h.userService.LoginByUsername(req)
	if err != nil {
		c.AbortWithStatusJSON(
			helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(token, true, helper.Success))
}

// RegisterByUsername
// @Summary      Register new user (username/email)
// @Description  Creates a new user account and sends a verification e-mail (if provided).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterByUsernameRequest  true  "Registration payload"
// @Success      201      {object}  helper.BaseHttpResponse                     "User created"
// @Failure      409      {object}  helper.BaseHttpResponse                     "Username or e-mail already exists"
// @Failure      422      {object}  helper.BaseHttpResponse                     "Validation error"
// @Failure      500      {object}  helper.BaseHttpResponse                     "Internal server error"
// @Router       /api/v1/users/register-by-username [post]
func (h *UserHandler) RegisterByUsername(c *gin.Context) {
	req := new(dto.RegisterByUsernameRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
		return
	}
	err = h.userService.RegisterByUsername(req)
	if err != nil {
		c.AbortWithStatusJSON(
			helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, helper.Success))
}

// RegisterLoginByMobileNumber
// @Summary      Register or login by mobile & OTP
// @Description  Verifies the OTP for the supplied mobile number.
//   - If the user exists ⇒ returns tokens.
//   - If the user does not exist ⇒ creates an account and returns tokens.
//
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterLoginByMobileRequest  true  "Mobile & OTP"
// @Success      201      {object}  helper.BaseHttpResponse{result=dto.TokenDetail}  "Tokens returned"
// @Failure      401      {object}  helper.BaseHttpResponse                            "Invalid or expired OTP"
// @Failure      422      {object}  helper.BaseHttpResponse                            "Validation error"
// @Failure      500      {object}  helper.BaseHttpResponse                            "Internal server error"
// @Router       /api/v1/users/login-by-mobile [post]
func (h *UserHandler) RegisterLoginByMobileNumber(c *gin.Context) {
	req := new(dto.RegisterLoginByMobileRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
		return
	}
	token, err := h.userService.RegisterLoginByMobileNumber(req)
	if err != nil {
		c.AbortWithStatusJSON(
			helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(token, true, helper.Success))
}

// Logout
// @Summary      Logout user
// @Description  Revokes both access and refresh tokens by blacklisting them in Redis until they expire.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LogoutRequest  true  "Refresh token"
// @Success      200      {object}  helper.BaseHttpResponse  "Logged out successfully"
// @Failure      401      {object}  helper.BaseHttpResponse  "Invalid or expired token"
// @Failure      422      {object}  helper.BaseHttpResponse  "Validation error"
// @Failure      500      {object}  helper.BaseHttpResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /api/v1/users/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	req := new(dto.LogoutRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err),
		)
		return
	}
	auth := c.GetHeader(constants.AuthorizationHeaderKey)
	accessToken, err := helper.ExtractToken(auth)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(nil, false, helper.AuthError,
				&service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}),
		)
		return
	}
	tokenSvc := services.NewTokenService(h.cfg)
	blackSvc := services.NewBlacklistService(h.cfg)
	acClaims, err := tokenSvc.GetClaims(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(nil, false, helper.AuthError, err))
		return
	}
	rtParsed, err := tokenSvc.VerifyRefreshToken(req.RefreshToken)
	if err != nil || !rtParsed.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			helper.GenerateBaseResponseWithError(nil, false, helper.AuthError,
				&service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}))
		return
	}
	rtClaims := rtParsed.Claims.(jwt.MapClaims)
	ttl := func(claims jwt.MapClaims) time.Duration {
		if exp, ok := claims[constants.ExpireTimeKey].(float64); ok {
			delta := time.Unix(int64(exp), 0).Sub(time.Now())
			if delta > 0 {
				return delta
			}
		}
		return 0
	}
	if err := blackSvc.Blacklist(accessToken, ttl(acClaims)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	if err := blackSvc.Blacklist(req.RefreshToken, ttl(rtClaims)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse("logged out", true, helper.Success))
}
