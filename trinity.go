package trinity

import (
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	// GlobalTrinity global instance
	GlobalTrinity  *Trinity
	runMode        = "Local"
	rootPath       string
	configFilePath string
)

func init() {
	rootPath, _ = os.Getwd()
	configFilePath = filepath.Join(rootPath, "config", "config.yml")
}

// Trinity struct for app subconfig
type Trinity struct {
	sync.RWMutex
	runMode        string
	router         *gin.Engine
	setting        *Setting
	db             *gorm.DB
	vCfg           *ViewSetCfg
	rootpath       string
	configFilePath string
	logger         Logger
}

func (t *Trinity) initDefaultValue() {
	GlobalTrinity = t
}

// GetRootPath  get rootpath
func GetRootPath() string {
	return rootPath
}

// SetRootPath  get rootpath
func SetRootPath(path string) {
	rootPath = path
}

// GetConfigFilePath  get rootpath
func GetConfigFilePath() string {
	return configFilePath
}

// SetConfigFilePath  get rootpath
func SetConfigFilePath(path string) {
	configFilePath = path
}

// GetRunMode  get RunMode
func GetRunMode() string {
	return runMode
}

// SetRunMode  set RunMode
func SetRunMode(runmode string) {
	runMode = runmode
}

// New app
// initial global trinity object
func New(customizeSettingSlice ...CustomizeSetting) *Trinity {
	t := &Trinity{
		runMode:        runMode,
		rootpath:       rootPath,
		configFilePath: configFilePath,
	}
	t.Lock()
	t.loadSetting(customizeSettingSlice...)
	t.initLogger()
	t.InitDatabase()
	t.initRouter()
	t.initViewSetCfg()
	t.initDefaultValue()
	t.Unlock()
	return t
}

// Reload  reload trinity
func (t *Trinity) Reload(runMode string) {
	t.RLock()
	t.runMode = runMode
	t.loadSetting()
	t.initLogger()
	t.InitDatabase()
	t.initRouter()
	t.initViewSetCfg()
	t.initDefaultValue()
	t.RUnlock()
}

// reloadTrinity for reload some config
func (t *Trinity) reloadTrinity() {
	t.loadSetting()
	t.initLogger()
	t.InitDatabase()
	t.initRouter()
	t.initViewSetCfg()
	t.initDefaultValue()
}

// GetVCfg  get vcfg
func (t *Trinity) GetVCfg() *ViewSetCfg {
	t.RLock()
	v := t.vCfg
	t.RUnlock()
	return v
}

// SetVCfg  get vcfg
func (t *Trinity) SetVCfg(newVCfg *ViewSetCfg) *Trinity {
	t.Lock()
	t.vCfg = newVCfg
	t.reloadTrinity()
	t.Unlock()
	return t
}

// GetSetting  get setting
func (t *Trinity) GetSetting() *Setting {
	t.RLock()
	s := t.setting
	t.RUnlock()
	return s
}

// SetSetting  set setting
func (t *Trinity) SetSetting(s *Setting) *Trinity {
	t.RLock()
	t.setting = s
	t.reloadTrinity()
	t.RUnlock()
	return t
}

// GetRouter  get router
func (t *Trinity) GetRouter() *gin.Engine {
	t.RLock()
	r := t.router
	t.RUnlock()
	return r
}

// SetRouter  set router
func (t *Trinity) SetRouter(newRouter *gin.Engine) *Trinity {
	t.Lock()
	t.router = newRouter
	t.reloadTrinity()
	t.Unlock()
	return t
}

// GetDB  get db instance
func (t *Trinity) GetDB() *gorm.DB {
	t.RLock()
	d := t.db
	t.RUnlock()
	return d
}

// SetDB  set db instance
func (t *Trinity) SetDB(db *gorm.DB) *Trinity {
	t.Lock()
	t.db = db
	t.reloadTrinity()
	t.Unlock()
	return t
}

// GetConfigFilePath  get rootpath
func (t *Trinity) GetConfigFilePath() string {
	t.RLock()
	r := t.configFilePath
	t.RUnlock()
	return r
}

// SetConfigFilePath  get rootpath
func (t *Trinity) SetConfigFilePath(configFilePath string) *Trinity {
	t.Lock()
	t.configFilePath = configFilePath
	t.reloadTrinity()
	t.Unlock()
	return t
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
