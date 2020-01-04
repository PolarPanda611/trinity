package trinity

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

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
	// db log
	SQLFunc    string `json:"sqlfunc"`
	SQL        string `json:"sql"`
	EffectRows int    `json:"effectrows"`
}

// Logger to record log
type Logger interface {
	Print(v ...interface{})
}

// defaultLogger: default logger
type defaultLogger struct {
}

// defaultLogger: default logger
type defaultViewRuntimeLogger struct {
	ViewRuntime *ViewSetRunTime
}

//CustomizeLogFormatter for customize log format
func CustomizeLogFormatter(params LogFormatterParams) string {
	l := LogFormat{
		Timestamp:      params.TimeStamp.Format(time.RFC3339),
		Version:        GlobalTrinity.setting.Version,
		Message:        params.ErrorMessage,
		LoggerName:     "",
		ThreadName:     "",
		Level:          "",
		Hostname:       "hostname",
		ModuleName:     GlobalTrinity.setting.Project,
		TraceID:        params.TraceID,
		Latency:        params.Latency,
		ClientIP:       params.ClientIP,
		HTTPMethod:     params.Method,
		HTTPPath:       params.Path,
		HTTPStatusCode: params.StatusCode,
		BodySize:       params.BodySize,
		UID:            params.UserID,
		ErrorDetail:    params.ErrorDetail,
		SQLFunc:        params.SQLFunc,
		SQL:            params.SQL,
		EffectRows:     params.EffectRows,
	}
	b, _ := json.Marshal(l)
	return string(b)

}

// InitLogger initial logger
func (t *Trinity) initLogger() {
	gin.SetMode(t.setting.Log.GinMode)
	runmode := gin.Mode()
	if runmode == "release" {
		if !CheckFileIsExist(t.setting.Log.LogRootPath) {
			if err := os.MkdirAll(t.setting.Log.LogRootPath, 770); err != nil {
				log.Fatalln("create log root path error：", err)
			}
		}
		var gFile *os.File
		var err error
		logfile := GetLogFilePath(t.setting.Log.LogRootPath, t.setting.Log.LogName)
		if CheckFileIsExist(logfile) {
			gFile, err = os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				log.Fatalln("open log error：", err)
			}

		} else {
			gFile, err = os.Create(logfile)
			if err != nil {
				log.Fatalln("create log error：", err)
			}
		}
		gin.DefaultWriter = io.MultiWriter(gFile)

	} else {
		gin.DefaultWriter = io.MultiWriter(os.Stderr)
	}
	t.logger = &defaultLogger{}
}

// LogWriter log
func (l *defaultLogger) Print(v ...interface{}) {
	//customie logger
	fmt.Println(v...)

}

// DbLoggerFormatter format gorm db log
func DbLoggerFormatter(r *ViewSetRunTime, v ...interface{}) {
	dblogLevel, _ := v[0].(string)
	sqlfunc := ""
	effectRows := 0
	if dblogLevel == "sql" {
		sqlfunc, _ = v[1].(string)
		// sql, _ := v[3].(string)
		effectRows, _ = v[5].(int)
		// sqldata := v[4]
	}
	l := LogFormat{
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   GlobalTrinity.setting.Version,
		// Message:        params.ErrorMessage,
		LoggerName: "",
		ThreadName: "",
		Level:      "",
		Hostname:   "hostname",
		ModuleName: GlobalTrinity.setting.Project,
		TraceID:    r.TraceID,
		// Latency:        params.Latency,
		ClientIP:       r.Gcontext.ClientIP(),
		HTTPMethod:     r.Gcontext.Request.Method,
		HTTPPath:       r.Gcontext.Request.URL.RequestURI(),
		HTTPStatusCode: r.Gcontext.Writer.Status(),
		BodySize:       r.Gcontext.Writer.Size(),
		UID:            r.Gcontext.GetString("UserID"),
		ErrorDetail:    r.Gcontext.GetString("ErrorDetail"),
		SQLFunc:        sqlfunc,
		SQL:            fmt.Sprintln(gorm.LogFormatter(v...)),
		EffectRows:     effectRows,
	}
	b, _ := json.Marshal(l)
	fmt.Fprintln(gin.DefaultWriter, string(b))

}

// LogWriter log
func (l *defaultViewRuntimeLogger) Print(v ...interface{}) {
	DbLoggerFormatter(l.ViewRuntime, v...)
}

// LogPrint customize log
func LogPrint(words string) {
	fmt.Fprintln(gin.DefaultWriter, words)
}
