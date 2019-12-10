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
	// ProjectRootPath project root path
	ProjectRootPath string
	// DefaultRunMode will load runmode config
	DefaultRunMode = "local"

	// DefaultJwtexpirehour for jwt
	DefaultJwtexpirehour = 2
	// DefaultJwtheaderprefix for jwt
	DefaultJwtheaderprefix = "Mold"
	// DefaultSecretkey for jwt
	DefaultSecretkey = "234"
	// DefaultJwtissuer for jwt
	DefaultJwtissuer = "Mold"
	// DefaultAppVersion for logging middleware
	DefaultAppVersion = "v1.0"
	// DefaultProjectName for logging middleware
	DefaultProjectName = "trinity"
	//Db global db instance
	Db *gorm.DB
	// GlobalTrinity global instance
	GlobalTrinity *Trinity
	// AppSetting Get Global setting
	AppSetting *Setting
)

// Trinity struct for app subconfig
type Trinity struct {
	sync.RWMutex
	RunMode string
	Router  *gin.Engine
	Setting *Setting
	Db      *gorm.DB
	vCfg    *ViewSetCfg
	Logger  Logger
}

func (t *Trinity) initDefaultValue() {
	ProjectRootPath, _ = os.Getwd()
	DefaultJwtexpirehour = t.Setting.Security.Authentication.JwtExpireHour
	DefaultJwtheaderprefix = t.Setting.Security.Authentication.JwtHeaderPrefix
	// DefaultSecretkey for jwt
	DefaultSecretkey = t.Setting.Security.Authentication.SecretKey
	// DefaultJwtissuer for jwt
	DefaultJwtissuer = t.Setting.Security.Authentication.JwtIssuer
	// DefaultAppVersion for logging middleware
	DefaultAppVersion = t.Setting.Version
	// DefaultProjectName for logging middleware
	DefaultProjectName = t.Setting.Project
	Db = t.Db
	AppSetting = t.Setting
	DefaultRunMode = t.RunMode
	GlobalTrinity = t
}

// New app
func New(runMode string) *Trinity {
	t := &Trinity{
		RunMode: runMode,
	}
	t.Lock()
	t.LoadSetting()
	t.InitLogger()
	t.InitDatabase()
	t.initDefaultValue()
	t.InitRouter()
	t.InitViewSetCfg()
	t.migrate()
	t.Unlock()
	return t
}

// Migrate run migration mode
func Migrate(runMode string) {
	t := &Trinity{
		RunMode: runMode,
	}
	t.Lock()
	t.LoadSetting()
	t.InitLogger()
	t.InitDatabase()
	t.initDefaultValue()
	t.migrate()
	t.Unlock()
	RunMigration()

}

// Serve http
func (t *Trinity) Serve() error {
	defer t.Close()
	s := &http.Server{
		Addr:              ":" + t.Setting.Webapp.Port,
		Handler:           t.Router,
		ReadTimeout:       time.Duration(t.Setting.Webapp.ReadTimeoutSecond) * time.Second,
		ReadHeaderTimeout: time.Duration(t.Setting.Webapp.ReadHeaderTimeoutSecond) * time.Second,
		WriteTimeout:      time.Duration(t.Setting.Webapp.WriteTimeoutSecond) * time.Second,
		IdleTimeout:       time.Duration(t.Setting.Webapp.IdleTimeoutSecond) * time.Second,
		MaxHeaderBytes:    t.Setting.Webapp.MaxHeaderBytes,
	}
	t.Logger.Print("[info]  " + time.Now().Format(time.RFC3339) + "  start http server listening : " + t.Setting.Webapp.Port + ", version : " + t.Setting.Version)
	return s.ListenAndServe()
}

// Close http
func (t *Trinity) Close() {
	t.Db.Close()
	t.Logger.Print("[info]  " + time.Now().Format(time.RFC3339) + "  end http server listening : " + t.Setting.Webapp.Port + ", version : " + t.Setting.Version)
}
