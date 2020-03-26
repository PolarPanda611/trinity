package trinity

import (
	"sync"

	"github.com/PolarPanda611/reflections"
	"github.com/jinzhu/configor"
)

// ISetting setting interface
type ISetting interface {
	GetDebug() bool
	GetTLSEnabled() bool
	GetSetting() *Setting
	GetProjectName() string
	GetTags() []string
	GetWebAppType() string
	GetWebAppAddress() string
	GetWebAppPort() int
	GetServiceMeshAddress() string
	GetServiceMeshPort() int
	GetDeregisterAfterCritical() int
	GetHealthCheckInterval() int
	GetCAPemFile() string
	GetServerPemFile() string
	GetServerKeyFile() string
	GetClientPemFile() string
	GetClientKeyFile() string
	GetProjectVersion() string
	GetLogRootPath() string
	GetLogName() string
	GetServiceMeshType() string
	GetServiceMeshAutoRegister() bool
	GetAtomicRequest() bool
	GetTablePrefix() string
	GetWebAppMediaURL() string
	GetWebAppStaticURL() string
	GetWebAppMediaPath() string
	GetWebAppStaticPath() string
	GetCacheSize() int
	GetCacheTimeout() int
	GetPageSize() int
	GetWebAppBaseURL() string
	GetMigrationPath() string
	GetJwtVerifyExpireHour() bool
	GetJwtVerifyIssuer() bool
	GetJwtIssuer() string
	GetJwtHeaderPrefix() string
	GetJwtExpireHour() int
	GetReadTimeoutSecond() int
	GetReadHeaderTimeoutSecond() int
	GetWriteTimeoutSecond() int
	GetIdleTimeoutSecond() int
	GetMaxHeaderBytes() int
	GetSecretKey() string
	GetAllowOrigins() []string
	GetAllowMethods() []string
	GetAllowHeaders() []string
	GetExposeHeaders() []string
	GetAllowCredentials() bool
	GetMaxAgeHour() int
	GetCorsEnable() bool
	GetDBHost() string
	GetDBPort() string
	GetDBUser() string
	GetDBPassword() string
	GetDBName() string
	GetDBOption() string
	GetDBType() string
	GetDbMaxIdleConn() int
	GetDbMaxOpenConn() int
}

// CustomizeSetting for customize setting
type CustomizeSetting interface {
	Load(runmode string, configFilePath string)
}

// Setting : for trinity setting
type Setting struct {
	mu      sync.RWMutex
	Project string   `yaml:"project"`
	Version string   `yaml:"version"`
	Tags    []string `yaml:"tags"`
	Runtime struct {
		Debug bool `yaml:"debug"`
	}
	Security struct {
		Authentication struct {
			SecretKey           string `yaml:"secretkey"`
			JwtVerifyIssuer     bool   `yaml:"jwtverifyissuer"`
			JwtIssuer           string `yaml:"jwtissuer"`
			JwtVerifyExpireHour bool   `yaml:"jwtverifyexpirehour"`
			JwtExpireHour       int    `yaml:"jwtexpirehour"`
			JwtHeaderPrefix     string `yaml:"jwtheaderprefix"`
		}
		Cors struct {
			Enable           bool     `yaml:"enable"`
			AllowOrigins     []string `yaml:"alloworigins"`
			AllowMethods     []string `yaml:"allowmethods"`
			AllowHeaders     []string `yaml:"allowheaders"`
			ExposeHeaders    []string `yaml:"exposeheaders"`
			AllowCredentials bool     `yaml:"allowcredentials"`
			MaxAgeHour       int      `yaml:"maxagehour"`
		}
		TLS struct {
			Enabled       bool   `yaml:"enabled"`
			CAPemFile     string `yaml:"ca_pem_file"`
			ServerPemFile string `yaml:"server_pem_file"`
			ServerKeyFile string `yaml:"server_key_file"`
			ClientPemFile string `yaml:"client_pem_file"`
			ClientKeyFile string `yaml:"client_key_file"`
		}
	}
	Webapp struct {
		// Type support GRPC HTTP
		Type    string `yaml:"type"`
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
		// ReadTimeout is the maximum duration for reading the entire
		// request, including the body.
		//
		// Because ReadTimeout does not let Handlers make per-request
		// decisions on each request body's acceptable deadline or
		// upload rate, most users will prefer to use
		// ReadHeaderTimeout. It is valid to use them both.
		ReadTimeoutSecond int `yaml:"readtimeoutsecond"`

		// ReadHeaderTimeout is the amount of time allowed to read
		// request headers. The connection's read deadline is reset
		// after reading the headers and the Handler can decide what
		// is considered too slow for the body. If ReadHeaderTimeout
		// is zero, the value of ReadTimeout is used. If both are
		// zero, there is no timeout.
		ReadHeaderTimeoutSecond int `yaml:"readheadertimeoutsecond"`

		// WriteTimeout is the maximum duration before timing out
		// writes of the response. It is reset whenever a new
		// request's header is read. Like ReadTimeout, it does not
		// let Handlers make decisions on a per-request basis.
		WriteTimeoutSecond int `yaml:"writertimeoutsecond"`

		// IdleTimeout is the maximum amount of time to wait for the
		// next request when keep-alives are enabled. If IdleTimeout
		// is zero, the value of ReadTimeout is used. If both are
		// zero, there is no timeout.
		IdleTimeoutSecond int `yaml:"idletimeoutsecond"`

		// MaxHeaderBytes controls the maximum number of bytes the
		// server will read parsing the request header's keys and
		// values, including the request line. It does not limit the
		// size of the request body.
		// If zero, DefaultMaxHeaderBytes is used.
		MaxHeaderBytes int    `yaml:"maxheaderbytes"`
		TemplatePath   string `yaml:"templatepath"`
		MediaURL       string `yaml:"mediaurl"`
		MediaPath      string `yaml:"mediapath"`
		StaticURL      string `yaml:"staticurl"`
		StaticPath     string `yaml:"staticpath"`
		MigrationPath  string `yaml:"migrationpath"`
		PageSize       int    `yaml:"pagesize"`
		MaxBodySize    int    `yaml:"maxbodysize"`
		AtomicRequest  bool   `yaml:"atomicrequest"`
		// if api root is not root , replease with base url
		// e.g : /assetgo
		BaseURL string `yaml:"baseurl"`
	}
	Log struct {
		LogRootPath string `yaml:"logrootpath"` //   /var/log/mold
		LogName     string `yaml:"logname"`     //  app.log
	}
	Cache struct {
		Redis struct {
			Host        string
			Port        int
			Password    string
			Maxidle     int
			Maxactive   int
			Idletimeout int
		}
		Gcache struct {
			CacheAlgorithm string `yaml:"cache_algorithm"`
			CacheSize      int    `yaml:"cachesize"`
			Timeout        int    `yaml:"timeout"` // hour
		}
	}
	Database struct {
		Type          string
		Name          string
		User          string
		Password      string
		Host          string
		Port          string
		Option        string
		TablePrefix   string
		DbMaxIdleConn int
		DbMaxOpenConn int
	}
	ServiceMesh struct {
		Type                    string // etcd oor consul
		Address                 string
		Port                    int
		DeregisterAfterCritical int  `yaml:"deregister_after_critical"` //second
		HealthCheckInterval     int  `yaml:"health_check_interval"`     //second
		AutoRegister            bool `yaml:"auto_register"`
	}
}

func (s *Setting) GetDbMaxIdleConn() int { return s.Database.DbMaxIdleConn }
func (s *Setting) GetDbMaxOpenConn() int { return s.Database.DbMaxOpenConn }
func (s *Setting) GetDBHost() string     { return s.Database.Host }
func (s *Setting) GetDBPort() string     { return s.Database.Port }
func (s *Setting) GetDBUser() string     { return s.Database.User }
func (s *Setting) GetDBPassword() string { return s.Database.Password }
func (s *Setting) GetDBName() string     { return s.Database.Name }
func (s *Setting) GetDBOption() string   { return s.Database.Option }
func (s *Setting) GetDBType() string     { return s.Database.Type }

func (s *Setting) GetCorsEnable() bool        { return s.Security.Cors.Enable }
func (s *Setting) GetMaxAgeHour() int         { return s.Security.Cors.MaxAgeHour }
func (s *Setting) GetAllowOrigins() []string  { return s.Security.Cors.AllowOrigins }
func (s *Setting) GetAllowMethods() []string  { return s.Security.Cors.AllowMethods }
func (s *Setting) GetAllowHeaders() []string  { return s.Security.Cors.AllowHeaders }
func (s *Setting) GetExposeHeaders() []string { return s.Security.Cors.ExposeHeaders }
func (s *Setting) GetAllowCredentials() bool {
	return s.Security.Cors.AllowCredentials
}

func (s *Setting) GetReadTimeoutSecond() int       { return s.Webapp.ReadTimeoutSecond }
func (s *Setting) GetReadHeaderTimeoutSecond() int { return s.Webapp.ReadHeaderTimeoutSecond }
func (s *Setting) GetWriteTimeoutSecond() int      { return s.Webapp.WriteTimeoutSecond }
func (s *Setting) GetIdleTimeoutSecond() int       { return s.Webapp.IdleTimeoutSecond }
func (s *Setting) GetMaxHeaderBytes() int          { return s.Webapp.MaxHeaderBytes }
func (s *Setting) GetSecretKey() string {
	return s.Security.Authentication.SecretKey
}

func (s *Setting) GetJwtExpireHour() int {
	return s.Security.Authentication.JwtExpireHour
}
func (s *Setting) GetJwtHeaderPrefix() string {
	return s.Security.Authentication.JwtHeaderPrefix
}
func (s *Setting) GetJwtIssuer() string {
	return s.Security.Authentication.JwtIssuer
}
func (s *Setting) GetJwtVerifyIssuer() bool {
	return s.Security.Authentication.JwtVerifyIssuer
}
func (s *Setting) GetJwtVerifyExpireHour() bool {
	return s.Security.Authentication.JwtVerifyExpireHour
}
func (s *Setting) GetMigrationPath() string { return s.Webapp.MigrationPath }
func (s *Setting) GetWebAppBaseURL() string { return s.Webapp.BaseURL }
func (s *Setting) GetPageSize() int         { return s.Webapp.PageSize }
func (s *Setting) GetCacheSize() int        { return s.Cache.Gcache.CacheSize }
func (s *Setting) GetCacheTimeout() int     { return s.Cache.Gcache.Timeout }

// GetWebAppMediaURL get web app media url
func (s *Setting) GetWebAppMediaURL() string { return s.Webapp.MediaURL }

// GetWebAppMediaPath get web app media path
func (s *Setting) GetWebAppMediaPath() string { return s.Webapp.MediaPath }

// GetWebAppStaticPath get web app static path
func (s *Setting) GetWebAppStaticPath() string { return s.Webapp.StaticPath }

// GetWebAppStaticURL get web app static url
func (s *Setting) GetWebAppStaticURL() string { return s.Webapp.StaticURL }

// GetLogRootPath get log root path
func (s *Setting) GetLogRootPath() string {
	return s.Log.LogRootPath
}

// GetTablePrefix get table prefix
func (s *Setting) GetTablePrefix() string {
	return s.Database.TablePrefix
}

// GetServiceMeshAutoRegister get auto register
func (s *Setting) GetServiceMeshAutoRegister() bool {
	return s.ServiceMesh.AutoRegister
}

// GetAtomicRequest get automic request is open
func (s *Setting) GetAtomicRequest() bool {
	return s.Webapp.AtomicRequest
}

//GetServiceMeshType get s m type
func (s *Setting) GetServiceMeshType() string {
	return s.ServiceMesh.Type
}

// GetTLSEnabled get tls enabled
func (s *Setting) GetTLSEnabled() bool {
	return s.Security.TLS.Enabled
}

// GetLogName get log name
func (s *Setting) GetLogName() string {
	return s.Log.LogName
}

// GetDebug get debug
func (s *Setting) GetDebug() bool {
	return s.Runtime.Debug
}

// GetSetting get setting
func (s *Setting) GetSetting() *Setting {
	return s
}

//GetCAPemFile get ca pem file
func (s *Setting) GetCAPemFile() string {
	return s.Security.TLS.CAPemFile
}

//GetServerPemFile get server pem file
func (s *Setting) GetServerPemFile() string {
	return s.Security.TLS.ServerPemFile
}

//GetServerKeyFile get server key file
func (s *Setting) GetServerKeyFile() string {
	return s.Security.TLS.ServerKeyFile
}

//GetClientPemFile get client pem file
func (s *Setting) GetClientPemFile() string {
	return s.Security.TLS.ClientPemFile
}

//GetClientKeyFile get client key file
func (s *Setting) GetClientKeyFile() string {
	return s.Security.TLS.ClientKeyFile
}

// GetDeregisterAfterCritical deregister service after critical second
func (s *Setting) GetDeregisterAfterCritical() int {
	return s.ServiceMesh.DeregisterAfterCritical
}

// GetHealthCheckInterval health check interval
func (s *Setting) GetHealthCheckInterval() int {
	return s.ServiceMesh.HealthCheckInterval

}

//GetTags get project tags
func (s *Setting) GetTags() []string {
	return s.Tags
}

// GetProjectName get project name
func (s *Setting) GetProjectName() string {
	return s.Project
}

// GetProjectVersion get project name
func (s *Setting) GetProjectVersion() string {
	return s.Version
}

// GetWebAppType get web app type
func (s *Setting) GetWebAppType() string {
	return s.Webapp.Type
}

// GetWebAppAddress get web service  ip address
func (s *Setting) GetWebAppAddress() string {
	return s.Webapp.Address
}

// GetWebAppPort get web service port
func (s *Setting) GetWebAppPort() int {
	return s.Webapp.Port
}

// GetServiceMeshAddress get service mesh address
func (s *Setting) GetServiceMeshAddress() string {
	return s.ServiceMesh.Address
}

// GetServiceMeshPort get service mesh port
func (s *Setting) GetServiceMeshPort() int {
	return s.ServiceMesh.Port
}

// GlobalSetting : for trinity global setting
type GlobalSetting struct {
	Local   Setting
	Develop Setting
	Preprod Setting
	Master  Setting
}

func newSetting(runMode string, configFilePath string) ISetting {
	s := GlobalSetting{}
	s.loadConfigFile(configFilePath)
	return s.loadSetting(runMode)
}

// loadConfigFile load config file
func (s *GlobalSetting) loadConfigFile(configFilePath string) {
	err := configor.Load(s, configFilePath)
	if err != nil {
		LoadConfigError(err)
	}
}

// loadSetting load config file
func (s *GlobalSetting) loadSetting(runMode string) ISetting {
	currentSettingInterface, err := reflections.GetField(s, runMode)
	if err != nil {
		WrongRunMode(runMode)
	}
	currentSetting, _ := currentSettingInterface.(Setting)
	return &currentSetting
}

// loadSetting used for load trinity config file by default and customize setting if necessery
func (t *Trinity) loadCustomizeSetting(customizeSettingSlice ...CustomizeSetting) {

	// load customize setting for application
	for _, v := range customizeSettingSlice {
		v.Load(t.runMode, t.configFilePath)
	}
}

// GetConfigFilePath  get rootpath
func (t *Trinity) GetConfigFilePath() string {
	t.mu.RLock()
	r := t.configFilePath
	t.mu.RUnlock()
	return r
}

// GetCurrentSetting  get setting
func (t *Trinity) GetCurrentSetting() ISetting {
	t.mu.RLock()
	s := t.setting
	t.mu.RUnlock()
	return s
}

// GetSetting  get setting
func (t *Trinity) GetSetting() ISetting {
	t.mu.RLock()
	s := t.setting
	t.mu.RUnlock()
	return s
}

// SetSetting  set setting
func (t *Trinity) SetSetting(s ISetting) *Trinity {
	t.mu.Lock()
	t.setting = s
	t.reloadTrinity()
	t.mu.Unlock()
	return t
}
