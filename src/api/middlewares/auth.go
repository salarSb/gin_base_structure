package middlewares

import (
	"base_structure/src/api/helper"
	"base_structure/src/config"
	"base_structure/src/constants"
	"base_structure/src/pkg/service_errors"
	"base_structure/src/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	tokenSvc := services.NewTokenService(cfg)
	blackSvc := services.NewBlacklistService(cfg)
	return func(c *gin.Context) {
		auth := c.GetHeader(constants.AuthorizationHeaderKey)
		if auth == "" {
			abortAuth(c, &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired})
			return
		}
		rawToken, err := helper.ExtractToken(auth)
		if err != nil {
			abortAuth(c, &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid})
			return
		}
		claims, err := tokenSvc.GetClaims(rawToken)
		if err != nil {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) && ve.Errors == jwt.ValidationErrorExpired {
				abortAuth(c, &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired})
			}
			return
		}
		if black, _ := blackSvc.IsBlacklisted(rawToken); black {
			abortAuth(c, &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid})
			return
		}
		for k, v := range claims {
			c.Set(k, v)
		}
		c.Next()
	}
}

func abortAuth(c *gin.Context, err error) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		helper.GenerateBaseResponseWithError(nil, false, helper.AuthError, err),
	)
}

func Authorization(validRoles []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		if len(context.Keys) == 0 {
			context.AbortWithStatusJSON(
				http.StatusForbidden,
				helper.GenerateBaseResponse(nil, false, helper.ForbiddenError),
			)
			return
		}
		rolesVal := context.Keys[constants.RolesKey]
		if rolesVal == nil {
			context.AbortWithStatusJSON(
				http.StatusForbidden,
				helper.GenerateBaseResponse(nil, false, helper.ForbiddenError),
			)
			return
		}
		roles := rolesVal.([]interface{})
		val := map[string]int{}
		for _, item := range roles {
			val[item.(string)] = 0
		}
		for _, item := range validRoles {
			if _, ok := val[item]; ok {
				context.Next()
				return
			}
		}
		context.AbortWithStatusJSON(
			http.StatusForbidden,
			helper.GenerateBaseResponse(nil, false, helper.ForbiddenError),
		)
	}
}
