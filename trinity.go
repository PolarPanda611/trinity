package trinity

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	// GlobalTrinity global instance
	GlobalTrinity *Trinity
)

// Trinity struct for app subconfig
type Trinity struct {
	sync.RWMutex
	runMode  string
	router   *gin.Engine
	setting  *Setting
	db       *gorm.DB
	vCfg     *ViewSetCfg
	rootPath string
	logger   Logger
}

func (t *Trinity) initDefaultValue() {
	GlobalTrinity = t
}

// New app
// initial global trinity object
func New(runMode string) *Trinity {
	rootPath, _ := os.Getwd()
	t := &Trinity{
		runMode:  runMode,
		rootPath: rootPath,
	}
	t.LoadSetting()
	t.InitLogger()
	t.InitDatabase()
	t.InitRouter()
	t.InitViewSetCfg()
	t.initDefaultValue()
	return t
}

// GetVCfg  get vcfg
func (t *Trinity) GetVCfg() *ViewSetCfg {
	t.RLock()
	v := t.vCfg
	t.RUnlock()
	return v
}

// SetVCfg  get vcfg
func (t *Trinity) SetVCfg(newVCfg *ViewSetCfg) {
	t.Lock()
	t.vCfg = newVCfg
	t.Unlock()
	return
}

// GetRunMode  get RunMode
func (t *Trinity) GetRunMode() string {
	t.RLock()
	r := t.runMode
	t.RUnlock()
	return r
}

// GetSetting  get setting
func (t *Trinity) GetSetting() *Setting {
	t.RLock()
	s := t.setting
	t.RUnlock()
	return s
}

// SetSetting  get setting
func (t *Trinity) SetSetting(s *Setting) {
	t.RLock()
	t.setting = s
	t.RUnlock()
}

// GetRouter  get router
func (t *Trinity) GetRouter() *gin.Engine {
	t.RLock()
	r := t.router
	t.RUnlock()
	return r
}

// SetRouter  set router
func (t *Trinity) SetRouter(newRouter *gin.Engine) {
	t.Lock()
	t.router = newRouter
	t.Unlock()
	return
}

// GetDB  get db instance
func (t *Trinity) GetDB() *gorm.DB {
	t.RLock()
	d := t.db
	t.RUnlock()
	return d
}

// SetDB  set db instance
func (t *Trinity) SetDB(db *gorm.DB) {
	t.Lock()
	t.db = db
	t.Unlock()
	return
}

// Migrate run migration mode
func Migrate(runMode string) {
	New(runMode)
	return

}

// Serve http
func (t *Trinity) Serve() error {
	defer t.Close()
	s := &http.Server{
		Addr:              ":" + t.setting.Webapp.Port,
		Handler:           t.router,
		ReadTimeout:       time.Duration(t.setting.Webapp.ReadTimeoutSecond) * time.Second,
		ReadHeaderTimeout: time.Duration(t.setting.Webapp.ReadHeaderTimeoutSecond) * time.Second,
		WriteTimeout:      time.Duration(t.setting.Webapp.WriteTimeoutSecond) * time.Second,
		IdleTimeout:       time.Duration(t.setting.Webapp.IdleTimeoutSecond) * time.Second,
		MaxHeaderBytes:    t.setting.Webapp.MaxHeaderBytes,
	}
	t.logger.Print("[info]  " + time.Now().Format(time.RFC3339) + "  start http server listening : " + t.setting.Webapp.Port + ", version : " + t.setting.Version)
	return s.ListenAndServe()
}

// Close http
func (t *Trinity) Close() {
	t.db.Close()
	t.logger.Print("[info]  " + time.Now().Format(time.RFC3339) + "  end http server listening : " + t.setting.Webapp.Port + ", version : " + t.setting.Version)
}
