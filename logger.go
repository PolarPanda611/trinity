package trinity

import "time"

import "fmt"

var (
	//DefaultLogger set default logger
	DefaultLogger Logger = &defaultLogger{}
)

// LogFormat log format struct
type LogFormat struct {
	Timestamp      string        `json:"@timestamp"` // current timestamp
	Version        string        `json:"@version"`   //
	LoggerName     string        `json:"logger_name"`
	ThreadName     string        `json:"thread_name"`
	Level          string        `json:"level"`
	Hostname       string        `json:"hostname"`
	ModuleName     string        `json:"module_name"`
	TraceID        string        `json:"trace_id"`
	Latency        time.Duration `json:"latency"`
	ClientIP       string        `json:"client_ip"`
	HTTPMethod     string        `json:"http_method"`
	HTTPPath       string        `json:"http_path"`
	HTTPStatusCode int           `json:"http_status_code"`
	Message        string        `json:"message"`      // error message
	ErrorDetail    string        `json:"error_detail"` // error detail info
	BodySize       int           `json:"body_size"`
	UID            string        `json:"uid"`
}

// Logger to record log
type Logger interface {
	Print(v ...interface{})
}

// defaultLogger: default logger
type defaultLogger struct{}

// LogWriter log
func (l *defaultLogger) Print(v ...interface{}) {
	//customie logger
	for _, m := range v {
		fmt.Println(m)
	}
}
