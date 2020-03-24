package trinity

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Logger to record log
type Logger interface {
	Print(v ...interface{})
}

// defaultLogger: default logger
type defaultLogger struct {
}

// LogWriter log
func (l *defaultLogger) Print(v ...interface{}) {
	LogPrint(v...)
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

// InitLogger initial logger
func (t *Trinity) initLogger() {
	if t.setting.Runtime.Debug {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}
	runmode := gin.Mode()
	if runmode == "release" {
		if !CheckFileIsExist(t.setting.Log.LogRootPath) {
			if err := os.MkdirAll(t.setting.Log.LogRootPath, 770); err != nil {
				log.Fatalln("create log root path error：", err)
			}
		}
		var gFile *os.File
		var err error
		logfile := filepath.Join(t.setting.Log.LogRootPath, t.setting.Log.LogName)
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

// LogPrint customize log
func LogPrint(v ...interface{}) {
	fmt.Fprintln(gin.DefaultWriter, v...)
}
