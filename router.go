package trinity

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// initRouter initial router
func (t *Trinity) initRouter() {
	// Creates a router without any middleware by default
	r := gin.New()
	// r.Use(timeoutMiddleware(time.Second * 2))
	r.Use(LogMiddleware())
	if t.setting.Security.Cors.Enable {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     t.setting.Security.Cors.AllowOrigins,
			AllowMethods:     t.setting.Security.Cors.AllowMethods,
			AllowHeaders:     t.setting.Security.Cors.AllowHeaders,
			ExposeHeaders:    t.setting.Security.Cors.ExposeHeaders,
			AllowCredentials: t.setting.Security.Cors.AllowCredentials,
			AllowOriginFunc: func(origin string) bool {
				return origin == "http://github.com"
			},
			MaxAge: time.Duration(t.setting.Security.Cors.MaxAgeHour) * time.Hour,
		}))
	}
	// r.LoadHTMLGlob(t.setting.Webapp.TemplatePath)
	r.RedirectTrailingSlash = false
	r.Use(gin.Recovery())
	r.Static(t.setting.Webapp.BaseURL+t.setting.Webapp.MediaURL, t.setting.Webapp.MediaPath)
	r.Static(t.setting.Webapp.BaseURL+t.setting.Webapp.StaticURL, t.setting.Webapp.StaticPath)
	r.GET(t.setting.Webapp.BaseURL+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET(t.setting.Webapp.BaseURL+"/api/ping", func(c *gin.Context) {
		err := GlobalTrinity.db.DB().Ping()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"Status":     400,
				"APIStatus":  "Error",
				"DBStatus":   "Error",
				"DBError":    err.Error(),
				"APIVersion": t.setting.Version,
			})
		} else {
			c.JSON(200, gin.H{
				"Status":     200,
				"APIStatus":  "alive",
				"DBStatus":   "alive",
				"DBInfo":     GlobalTrinity.db.DB().Stats(),
				"APIVersion": t.setting.Version,
			})
		}
	})
	t.router = r
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

// NewAPIGroup register new apigroup
func (t *Trinity) NewAPIGroup(path string) *gin.RouterGroup {
	return t.router.Group(t.setting.Webapp.BaseURL + path)
}

// NewAPIInGroup register new api in group
func NewAPIInGroup(rg *gin.RouterGroup, resource string, viewset gin.HandlerFunc, SupportedMethod []string) {
	SupportMethodMap := map[string]bool{
		"Retrieve": false,
		"List":     false,
		"Create":   false,
		"Update":   false,
		"Delete":   false,
	}
	for _, v := range SupportedMethod {
		_, exist := SupportMethodMap[v]
		if exist {
			SupportMethodMap[v] = true
		}
	}
	if SupportMethodMap["Retrieve"] {
		rg.GET("/"+resource+"/:id", viewset)
	}
	if SupportMethodMap["List"] {
		rg.GET("/"+resource, viewset)
	}
	if SupportMethodMap["Create"] {
		rg.POST("/"+resource, viewset)
	}
	if SupportMethodMap["Update"] {
		rg.PATCH("/"+resource+"/:id", viewset)
		rg.PUT("/"+resource+"/:id", viewset)
	}
	if SupportMethodMap["Delete"] {
		rg.DELETE("/"+resource+"/:id", viewset)
	}
}
