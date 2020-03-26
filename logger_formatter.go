package trinity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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
	UserID         int64         `json:"user_id,string"`
	Username       string        `json:"username"`
	// db log
	DBRunningFile  string        `json:"db_running_file"`
	DBRunningTime  time.Duration `json:"db_running_time"`
	DBSQL          string        `json:"db_sql"`
	DBParams       string        `json:"db_params"`
	DBEffectedRows string        `json:"db_effected_rows"`
	DBLogOrigin    string        `json:"db_log_origin"`
}

// GetString get logformat to string
func (l *LogFormat) GetString() string {
	b, _ := json.Marshal(l)
	return string(b)
}

// DbLoggerFormatter  for db execute logger
func DbLoggerFormatter(r *ViewSetRunTime, v ...interface{}) string {

	dblogLevel, _ := v[0].(string)
	l := &LogFormat{
		Timestamp:  time.Now().Format(time.RFC3339),
		Version:    GlobalTrinity.setting.GetProjectVersion(),
		Message:    r.Gcontext.Errors.String(),
		LoggerName: "",
		ThreadName: "",
		Level:      "",
		Hostname:   "hostname",
		ModuleName: GlobalTrinity.setting.GetProjectName(),
		TraceID:    r.TraceID,
		// Latency:        params.Latency,
		ClientIP:       r.Gcontext.ClientIP(),
		HTTPMethod:     r.Gcontext.Request.Method,
		HTTPPath:       r.Gcontext.Request.URL.RequestURI(),
		HTTPStatusCode: r.Gcontext.Writer.Status(),
		BodySize:       r.Gcontext.Writer.Size(),
		UserID:         r.Gcontext.GetInt64("UserID"),
		Username:       r.Gcontext.GetString("Username"),
		ErrorDetail:    r.Gcontext.GetString("ErrorDetail"),
		DBLogOrigin:    fmt.Sprint(gorm.LogFormatter(v...)),
	}
	if dblogLevel == "sql" {
		l.DBRunningFile = fmt.Sprint(v[1])
		l.DBRunningTime, _ = v[2].(time.Duration)
		l.DBSQL = fmt.Sprint(v[3])
		l.DBParams = fmt.Sprint(v[4])
		l.DBEffectedRows = fmt.Sprint(v[5])
	}
	return l.GetString()
}
