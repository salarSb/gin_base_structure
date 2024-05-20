package middlewares

import (
	"base_structure/src/api/helper"
	"base_structure/src/config"
	"base_structure/src/pkg/limiter"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func OtpLimiter(cfg *config.Config) gin.HandlerFunc {
	var ipLimiter = limiter.NewIpRateLimiter(rate.Every(cfg.Otp.Limiter*time.Second), 1)
	return func(context *gin.Context) {
		ipLimiter := ipLimiter.GetLimiter(context.Request.RemoteAddr)
		if !ipLimiter.Allow() {
			context.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				helper.GenerateBaseResponseWithError(
					nil,
					false,
					helper.OtpLimiterError,
					errors.New("not allowed"),
				),
			)
			context.Abort()
		} else {
			context.Next()
		}
	}
}
