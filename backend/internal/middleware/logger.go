package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"timestamp":  param.TimeStamp.Format(time.RFC3339),
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"method":     param.Method,
			"path":       param.Path,
			"user_agent": param.Request.UserAgent(),
		}).Info("HTTP Request")

		return ""
	})
}
