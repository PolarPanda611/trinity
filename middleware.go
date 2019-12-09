package trinity

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// Error Detail
	ErrorDetail string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
	// UserID is the request user
	UserID string
	// request TraceID unique
	TraceID string

	// db log
	SQLFunc    string `json:"sqlfunc"`
	SQL        string `json:"sql"`
	EffectRows int    `json:"effectrows"`
}

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter

	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	isTerm := true

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		raw := c.Request.URL.RequestURI()
		c.Set("TraceID", uuid.NewV4().String())

		// Process request
		c.Next()

		param := LogFormatterParams{
			Request: c.Request,
			isTerm:  isTerm,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.String()

		param.BodySize = c.Writer.Size()
		param.UserID = c.GetString("UserID")
		param.TraceID = c.GetString("TraceID")
		param.Path = raw
		param.ErrorDetail = c.GetString("ErrorDetail")
		fmt.Fprintln(out, formatter(param))

	}
}

// LoggerWithFormatter instance a Logger middleware with the specified log format function.
func LoggerWithFormatter() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Formatter: CustomizeLogFormatter,
	})
}

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := CheckTokenValid(c)

		if err != nil {
			c.AbortWithError(401, err)
			return
		}

		c.Next()
	}
}
