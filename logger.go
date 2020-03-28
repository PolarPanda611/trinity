package trinity

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/kit/log"
)

// Logger to record log
type Logger interface {
	Print(v ...interface{})
}

// NilLogger nil logger
type NilLogger struct{}

// Print nil logger do noothing
func (l *NilLogger) Print(v ...interface{}) {

}

// defaultLogger: default logger
type defaultLogger struct {
	userRequestsCtx UserRequestsCtx
	setting         ISetting
}

// NewDefaultLogger new default logger
func NewDefaultLogger(userRequestsCtx UserRequestsCtx, setting ISetting) Logger {
	return &defaultLogger{
		userRequestsCtx: userRequestsCtx,
		setting:         setting,
	}

}

// LogWriter log
func (l *defaultLogger) Print(v ...interface{}) {

	var logInterface []interface{}
	logInterface = []interface{}{
		"ServiceName=", GetServiceName(l.setting.GetProjectName()),
		"Time=", kitlog.DefaultTimestamp(),
		"Caller=", kitlog.DefaultCaller(),
		"Method=", l.userRequestsCtx.GetGRPCMethod(),
		"TraceID=", l.userRequestsCtx.GetTraceID(),
		"ReqUserName=", l.userRequestsCtx.GetReqUserName(),
	}
	if len(v) > 0 {
		dblogLevel, _ := v[0].(string)
		if dblogLevel == "sql" {
			logInterface = append(logInterface, "DBRunningFile=")
			logInterface = append(logInterface, fmt.Sprint(v[1]))
			logInterface = append(logInterface, "DBRunningTime=")
			DBRunningTime, _ := v[2].(time.Duration)
			logInterface = append(logInterface, DBRunningTime)
			logInterface = append(logInterface, "DBSQL=")
			logInterface = append(logInterface, fmt.Sprint(v[3]))
			logInterface = append(logInterface, "DBParams=")
			logInterface = append(logInterface, fmt.Sprint(v[4]))
			logInterface = append(logInterface, "DBEffectedRows=")
			logInterface = append(logInterface, fmt.Sprint(v[5]))
		}
	}
	logInterface = append(logInterface, v...)
	// l.Logger.Log(logInterface...)
	// fmt.Fprintln(gin.DefaultWriter, logInterface...)
	LogPrint(logInterface...)
}

// defaultViewRuntimeLogger: default logger for request
type defaultViewRuntimeLogger struct {
	ViewRuntime *ViewSetRunTime
}

// LogWriter log
func (l *defaultViewRuntimeLogger) Print(v ...interface{}) {
	log := DbLoggerFormatter(l.ViewRuntime, v...)
	LogPrint(log)
}

// initLoggerWriter initial logger file
func initLoggerWriter(setting ISetting) io.Writer {
	if setting.GetDebug() {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}
	runmode := gin.Mode()
	if runmode == "release" {
		if !CheckFileIsExist(setting.GetLogRootPath()) {
			if err := os.MkdirAll(setting.GetLogRootPath(), 770); err != nil {
				log.Fatalln("create log root path error：", err)
			}
		}
		var gFile *os.File
		var err error
		logfile := filepath.Join(setting.GetLogRootPath(), setting.GetLogName())
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
		return io.MultiWriter(gFile)

	}
	return io.MultiWriter(os.Stderr)

}

// LogPrint customize log
func LogPrint(v ...interface{}) {
	fmt.Fprintln(defaultWriter, v...)
}
