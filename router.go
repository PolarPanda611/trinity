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
	if t.setting.GetCorsEnable() {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     t.setting.GetAllowOrigins(),
			AllowMethods:     t.setting.GetAllowMethods(),
			AllowHeaders:     t.setting.GetAllowHeaders(),
			ExposeHeaders:    t.setting.GetExposeHeaders(),
			AllowCredentials: t.setting.GetAllowCredentials(),
			AllowOriginFunc: func(origin string) bool {
				return origin == "http://github.com"
			},
			MaxAge: time.Duration(t.setting.GetMaxAgeHour()) * time.Hour,
		}))

	}
	// r.LoadHTMLGlob(t.setting.Webapp.TemplatePath)
	r.RedirectTrailingSlash = false
	r.Use(gin.Recovery())
	r.Static(t.setting.GetWebAppBaseURL()+t.setting.GetWebAppMediaURL(), t.setting.GetWebAppMediaPath())
	r.Static(t.setting.GetWebAppBaseURL()+t.setting.GetWebAppStaticURL(), t.setting.GetWebAppStaticPath())
	r.GET(t.setting.GetWebAppBaseURL()+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET(t.setting.GetWebAppBaseURL()+"/api/ping", func(c *gin.Context) {
		err := GlobalTrinity.db.DB().Ping()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"Status":     400,
				"APIStatus":  "Error",
				"DBStatus":   "Error",
				"DBError":    err.Error(),
				"APIVersion": t.setting.GetProjectVersion(),
			})
		} else {
			c.JSON(200, gin.H{
				"Status":     200,
				"APIStatus":  "alive",
				"DBStatus":   "alive",
				"DBInfo":     GlobalTrinity.db.DB().Stats(),
				"APIVersion": t.setting.GetProjectVersion(),
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

// NewAPIGroup register new apigroup
func (t *Trinity) NewAPIGroup(path string) *gin.RouterGroup {
	return t.router.Group(t.setting.GetWebAppBaseURL() + path)
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
