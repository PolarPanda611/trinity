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

var logger kitlog.Logger

// Logger to record log
type Logger interface {
	clone() Logger
	FormatLogger(method string, traceID string, user string) Logger
	Print(v ...interface{})
}

// defaultLogger: default logger
type defaultLogger struct {
	Out            io.Writer
	Logger         kitlog.Logger
	ProjectName    string
	ProjectVersion string
	WebAppAddress  string
	WebAppPort     int
	Method         string
	TraceID        string
	User           string
}

func (l *defaultLogger) clone() Logger {
	return &defaultLogger{
		Out:            l.Out,
		Logger:         kitlog.NewLogfmtLogger(l.Out),
		ProjectName:    l.ProjectName,
		ProjectVersion: l.ProjectVersion,
		WebAppAddress:  l.WebAppAddress,
		WebAppPort:     l.WebAppPort,
	}
}

func (l *defaultLogger) FormatLogger(method string, traceID string, user string) Logger {

	l.Method = method
	l.TraceID = traceID
	l.User = user
	// l.Logger = kitlog.With(l.Logger, "ServiceName", GetServiceName(l.ProjectName, l.ProjectVersion))
	// l.Logger = kitlog.With(l.Logger, "Time", kitlog.DefaultTimestampUTC)
	// l.Logger = kitlog.With(l.Logger, "Caller", kitlog.DefaultCaller)
	// l.Logger = kitlog.With(l.Logger, "Method", method)
	// l.Logger = kitlog.With(l.Logger, "TraceID", traceID)
	// l.Logger = kitlog.With(l.Logger, "User", user)
	return l
	// logger.Log(v...)
}

// LogWriter log
func (l *defaultLogger) Print(v ...interface{}) {
	var logInterface []interface{}
	logInterface = []interface{}{
		"ServiceName", GetServiceName(l.ProjectName, l.ProjectVersion),
		"Time", kitlog.DefaultTimestamp(),
		"Caller", kitlog.DefaultCaller(),
		"Method", l.Method,
		"TraceID", l.TraceID,
		"User", l.User,
	}
	if len(v) > 0 {
		dblogLevel, _ := v[0].(string)
		if dblogLevel == "sql" {
			logInterface = append(logInterface, "DBRunningFile")
			logInterface = append(logInterface, fmt.Sprint(v[1]))
			logInterface = append(logInterface, "DBRunningTime")
			DBRunningTime, _ := v[2].(time.Duration)
			logInterface = append(logInterface, DBRunningTime)
			logInterface = append(logInterface, "DBSQL")
			logInterface = append(logInterface, fmt.Sprint(v[3]))
			logInterface = append(logInterface, "DBParams")
			logInterface = append(logInterface, fmt.Sprint(v[4]))
			logInterface = append(logInterface, "DBEffectedRows")
			logInterface = append(logInterface, fmt.Sprint(v[5]))
		}
	}
	logInterface = append(logInterface, v...)
	l.Logger.Log(logInterface...)
	// LogPrint(v...)
}

// defaultViewRuntimeLogger: default logger for request
type defaultViewRuntimeLogger struct {
	ViewRuntime *ViewSetRunTime
}

func (l *defaultViewRuntimeLogger) clone() Logger {
	return l
}
func (l *defaultViewRuntimeLogger) FormatLogger(method string, traceID string, user string) Logger {
	return l
}

// LogWriter log
func (l *defaultViewRuntimeLogger) Print(v ...interface{}) {
	log := DbLoggerFormatter(l.ViewRuntime, v...)
	LogPrint(log)
}

// InitLogger initial logger
func initLogger(setting ISetting) Logger {
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
		gin.DefaultWriter = io.MultiWriter(gFile)

	} else {
		gin.DefaultWriter = io.MultiWriter(os.Stderr)
	}
	return &defaultLogger{
		Out:            kitlog.NewSyncWriter(gin.DefaultWriter),
		Logger:         kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(gin.DefaultWriter)),
		ProjectName:    setting.GetProjectName(),
		ProjectVersion: setting.GetProjectVersion(),
		WebAppAddress:  setting.GetWebAppAddress(),
		WebAppPort:     setting.GetWebAppPort(),
	}
}

// LogPrint customize log
func LogPrint(v ...interface{}) {
	fmt.Fprintln(gin.DefaultWriter, v...)
}
