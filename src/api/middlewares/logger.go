package middlewares

import (
	"base_structure/src/pkg/logging"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"mime"
	"time"
)

const (
	maxBodySize     = 32 * 1024 // 32 KB
	skipContentType = "multipart/form-data"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if err == nil && w.buf.Len() < maxBodySize {
		_, _ = w.buf.Write(b)
	}
	return n, err
}

func StructuredLogger(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/swagger/*any" {
			c.Next()
			return
		}
		var reqBody string
		if ct, _, _ := mime.ParseMediaType(c.ContentType()); ct != skipContentType {
			lr := io.LimitReader(c.Request.Body, maxBodySize+1)
			b, _ := io.ReadAll(lr)
			_ = c.Request.Body.Close()
			c.Request.Body = io.NopCloser(bytes.NewReader(b))
			if len(b) > maxBodySize {
				reqBody = "(truncated)"
			} else {
				reqBody = string(b)
			}
		}
		blw := &bodyLogWriter{ResponseWriter: c.Writer, buf: bytes.NewBuffer(make([]byte, 0, 256))}
		c.Writer = blw
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		fields := map[logging.ExtraKey]interface{}{
			logging.Path:        c.FullPath(),
			logging.ClientIp:    c.ClientIP(),
			logging.Method:      c.Request.Method,
			logging.StatusCode:  c.Writer.Status(),
			logging.Latency:     latency,
			logging.BodySize:    c.Writer.Size(),
			logging.RequestBody: reqBody,
		}
		if em := c.Errors.ByType(gin.ErrorTypePrivate).String(); em != "" {
			fields[logging.ErrorMessage] = em
		}
		if ct, _, _ := mime.ParseMediaType(c.Writer.Header().Get("Content-Type")); ct == "application/json" && blw.buf.Len() < maxBodySize {
			fields[logging.ResponseBody] = blw.buf.String()
		}
		logger.Info(logging.RequestResponse, logging.Api, "", fields)
	}
}
