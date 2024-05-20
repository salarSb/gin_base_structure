package middlewares

import (
	"base_structure/src/config"
	"base_structure/src/pkg/logging"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func DefaultStructuredLogger(cfg *config.Config) gin.HandlerFunc {
	logger := logging.NewLogger(cfg)
	return structuredLogger(logger)
}

func structuredLogger(logger logging.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		if strings.Contains(context.FullPath(), "swagger") {
			context.Next()
		} else {
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
			start := time.Now() //start
			path := context.FullPath()
			raw := context.Request.URL.RawQuery
			bodyBytes, _ := io.ReadAll(context.Request.Body)
			err := context.Request.Body.Close()
			if err != nil {
				logger.Fatal(
					logging.RequestResponse,
					logging.Api,
					"cannot close byte request body in structuredLogger middleware",
					nil,
				)
				return
			}
			context.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			context.Writer = blw
			context.Next()
			param := gin.LogFormatterParams{}
			param.TimeStamp = time.Now() // stop
			param.Latency = param.TimeStamp.Sub(start)
			param.ClientIP = context.ClientIP()
			param.Method = context.Request.Method
			param.StatusCode = context.Writer.Status()
			param.ErrorMessage = context.Errors.ByType(gin.ErrorTypePrivate).String()
			param.BodySize = context.Writer.Size()
			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path
			keys := map[logging.ExtraKey]interface{}{}
			keys[logging.Path] = param.Path
			keys[logging.ClientIp] = param.ClientIP
			keys[logging.Method] = param.Method
			keys[logging.Latency] = param.Latency
			keys[logging.StatusCode] = param.StatusCode
			keys[logging.ErrorMessage] = param.ErrorMessage
			keys[logging.BodySize] = param.BodySize
			keys[logging.RequestBody] = string(bodyBytes)
			keys[logging.ResponseBody] = blw.body.String()
			logger.Info(logging.RequestResponse, logging.Api, "", keys)
		}
	}
}
