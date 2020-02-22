package trinity

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bluele/gcache"
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
	mu             sync.RWMutex
	runMode        string
	router         *gin.Engine
	setting        *Setting
	db             *gorm.DB
	vCfg           *ViewSetCfg
	rootpath       string
	configFilePath string
	logger         Logger
	cache          gcache.Cache
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
	t.mu.Lock()
	t.loadSetting(customizeSettingSlice...)
	t.initLogger()
	t.InitDatabase()
	t.initRouter()
	t.initViewSetCfg()
	t.initDefaultValue()
	t.cache = gcache.New(t.setting.Cache.Gcache.CacheSize).LRU().Expiration(time.Duration(t.setting.Cache.Gcache.Timeout) * time.Hour).Build()
	t.mu.Unlock()
	return t
}

// Reload  reload trinity
func (t *Trinity) Reload(runMode string) {
	t.mu.RLock()
	t.runMode = runMode
	t.loadSetting()
	t.initLogger()
	t.InitDatabase()
	t.initRouter()
	t.initViewSetCfg()
	t.initDefaultValue()
	t.mu.RUnlock()
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
	t.mu.RLock()
	v := t.vCfg
	t.mu.RUnlock()
	return v
}

// GetCache  get vcfg
func (t *Trinity) GetCache() gcache.Cache {
	t.mu.RLock()
	c := t.cache
	t.mu.RUnlock()
	return c
}

// SetCache  get vcfg
func (t *Trinity) SetCache(cache gcache.Cache) *Trinity {
	t.mu.Lock()
	t.cache = cache
	t.mu.Unlock()
	return t
}

// SetVCfg  get vcfg
func (t *Trinity) SetVCfg(newVCfg *ViewSetCfg) *Trinity {
	t.mu.Lock()
	t.vCfg = newVCfg
	t.reloadTrinity()
	t.mu.Unlock()
	return t
}

// GetSetting  get setting
func (t *Trinity) GetSetting() *Setting {
	t.mu.RLock()
	s := t.setting
	t.mu.RUnlock()
	return s
}

// SetSetting  set setting
func (t *Trinity) SetSetting(s *Setting) *Trinity {
	t.mu.RLock()
	t.setting = s
	t.reloadTrinity()
	t.mu.RUnlock()
	return t
}

// GetRouter  get router
func (t *Trinity) GetRouter() *gin.Engine {
	t.mu.RLock()
	r := t.router
	t.mu.RUnlock()
	return r
}

// SetRouter  set router
func (t *Trinity) SetRouter(newRouter *gin.Engine) *Trinity {
	t.mu.Lock()
	t.router = newRouter
	t.reloadTrinity()
	t.mu.Unlock()
	return t
}

// GetDB  get db instance
func (t *Trinity) GetDB() *gorm.DB {
	t.mu.RLock()
	d := t.db
	t.mu.RUnlock()
	return d
}

// SetDB  set db instance
func (t *Trinity) SetDB(db *gorm.DB) *Trinity {
	t.mu.Lock()
	t.db = db
	t.reloadTrinity()
	t.mu.Unlock()
	return t
}

// GetConfigFilePath  get rootpath
func (t *Trinity) GetConfigFilePath() string {
	t.mu.RLock()
	r := t.configFilePath
	t.mu.RUnlock()
	return r
}

// SetConfigFilePath  get rootpath
func (t *Trinity) SetConfigFilePath(configFilePath string) *Trinity {
	t.mu.Lock()
	t.configFilePath = configFilePath
	t.reloadTrinity()
	t.mu.Unlock()
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
	t.logger.Print(fmt.Sprintf("[info] %v  start http server listening : %v, version : %v", GetCurrentTimeString(time.RFC3339), t.setting.Webapp.Port, t.setting.Version))
	return s.ListenAndServe()
}

// Close http
func (t *Trinity) Close() {
	t.db.Close()
	t.logger.Print(fmt.Sprintf("[info] %v  end http server listening : %v, version : %v", GetCurrentTimeString(time.RFC3339), t.setting.Webapp.Port, t.setting.Version))
}
