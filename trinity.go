package trinity

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
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
	DefaultRunMode = t.RunMode
	GlobalTrinity = t
}

// New app
func New(runMode string) *Trinity {
	logger := &defaultLogger{}
	t := &Trinity{
		RunMode: runMode,
		Logger:  logger,
	}
	t.Lock()
	t.LoadSetting()
	t.InitDatabase()
	t.initDefaultValue()
	t.NewRouter()
	t.InitViewSetCfg()
	t.migrate()
	v1 := t.NewAPIGroup("/api/v1")
	NewAPIInGroup(v1, "users", UserViewSet, []string{"Retrieve", "List", "Create", "Update", "Delete"})
	NewAPIInGroup(v1, "permissions", PermissionViewSet, []string{"Retrieve", "List", "Create", "Update", "Delete"})
	NewAPIInGroup(v1, "groups", GroupViewSet, []string{"Retrieve", "List", "Create", "Update", "Delete"})
	NewAPIInGroup(v1, "apperrors", AppErrorViewSet, []string{"Retrieve", "List"})
	t.Unlock()
	return t
}

// Serve http
func (t *Trinity) Serve() {
	s := &http.Server{
		Addr:              ":" + t.Setting.HTTP.Port,
		Handler:           t.Router,
		ReadTimeout:       time.Duration(t.Setting.HTTP.ReadTimeoutSecond) * time.Second,
		ReadHeaderTimeout: time.Duration(t.Setting.HTTP.ReadHeaderTimeoutSecond) * time.Second,
		WriteTimeout:      time.Duration(t.Setting.HTTP.WriteTimeoutSecond) * time.Second,
		IdleTimeout:       time.Duration(t.Setting.HTTP.IdleTimeoutSecond) * time.Second,
		MaxHeaderBytes:    t.Setting.HTTP.MaxHeaderBytes,
	}
	t.Logger.Print("[info] %v start http server listening : %v , version : %v ", time.Now().Format(time.RFC3339), t.Setting.HTTP.Port, t.Setting.Version)
	fmt.Printf("[info] %v start http server listening : %v , version : %v ", time.Now().Format(time.RFC3339), t.Setting.HTTP.Port, t.Setting.Version)
	s.ListenAndServe()
	return
}

// Close http
func (t *Trinity) Close() {
	fmt.Printf("[info] %v end http server listening : %v , version : %v ", time.Now().Format(time.RFC3339), t.Setting.HTTP.Port, t.Setting.Version)
	t.Logger.Print("[info] %v end http server listening : %v , version : %v ", time.Now().Format(time.RFC3339), t.Setting.HTTP.Port, t.Setting.Version)
	return
}
