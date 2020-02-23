package trinity

import (
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// LogMiddleware  gin log formatter
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := GetCurrentTime()
		raw := c.Request.URL.RequestURI()
		traceID := uuid.NewV4().String()
		c.Set("TraceID", traceID)
		// Process request
		c.Next()

		//stop timmer
		currentTime := GetCurrentTime()
		currentTimeString := currentTime.Format(time.RFC3339)
		l := &LogFormat{
			Timestamp:      currentTimeString,
			Version:        GlobalTrinity.setting.Version,
			Message:        c.Errors.String(),
			LoggerName:     "",
			ThreadName:     "",
			Level:          "",
			Hostname:       "hostname",
			ModuleName:     GlobalTrinity.setting.Project,
			TraceID:        traceID,
			Latency:        currentTime.Sub(start),
			ClientIP:       c.ClientIP(),
			HTTPMethod:     c.Request.Method,
			HTTPPath:       raw,
			HTTPStatusCode: c.Writer.Status(),
			BodySize:       c.Writer.Size(),
			UserID:         c.GetInt64("UserID"),
			Username:       c.GetString("Username"),
			ErrorDetail:    c.GetString("ErrorDetail"),
		}

		LogPrint(l.GetString())

	}
}
