package handlers

import (
	"base_structure/src/api/dto"
	"base_structure/src/api/helper"
	"base_structure/src/config"
	"base_structure/src/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	userService := services.NewUserService(cfg)
	return &UserHandler{userService: userService}
}

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
