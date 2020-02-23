package trinity

import (
	"github.com/PolarPanda611/reflections"
	"github.com/jinzhu/configor"
)

// CustomizeSetting for customize setting
type CustomizeSetting interface {
	Load(runmode string, configFilePath string)
}

// Setting : for trinity setting
type Setting struct {
	Project string `yaml:"project"`
	Version string `yaml:"version"`
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
	}
	Webapp struct {
		Port string `yaml:"port"`
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
		GinMode     string `yaml:"ginmode"`     //release   log=>log file , debug, test : log=>console
		LogRootPath string `yaml:"logrootpath"` //   /var/log/mold
		LogName     string `yaml:"logname"`     //  app.log
	}
	Cache struct {
		Redis struct {
			Host        string
			Port        string
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
}

// GlobalSetting : for trinity global setting
type GlobalSetting struct {
	Local   Setting
	Develop Setting
	Preprod Setting
	Master  Setting
}

// loadSetting used for load trinity config file by default and customize setting if necessery
func (t *Trinity) loadSetting(customizeSettingSlice ...CustomizeSetting) {

	// load global setting for trinity
	g := GlobalSetting{}
	err := configor.Load(&g, t.configFilePath)
	if err != nil {
		LoadConfigError(err)
	}
	currentSettingInterface, err := reflections.GetField(g, t.runMode)
	if err != nil {
		WrongRunMode(t.runMode)
	}
	currentSetting, _ := currentSettingInterface.(Setting)
	t.setting = &currentSetting

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

// SetConfigFilePath  get rootpath
func (t *Trinity) SetConfigFilePath(configFilePath string) *Trinity {
	t.mu.Lock()
	t.configFilePath = configFilePath
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
	t.mu.Lock()
	t.setting = s
	t.reloadTrinity()
	t.mu.Unlock()
	return t
}
