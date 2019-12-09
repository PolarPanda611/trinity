package trinity

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// NewRouter initial router
func (t *Trinity) NewRouter() {
	// Creates a router without any middleware by default
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     t.Setting.Security.Cors.AllowOrigins,
		AllowMethods:     t.Setting.Security.Cors.AllowMethods,
		AllowHeaders:     t.Setting.Security.Cors.AllowHeaders,
		ExposeHeaders:    t.Setting.Security.Cors.ExposeHeaders,
		AllowCredentials: t.Setting.Security.Cors.AllowCredentials,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: time.Duration(t.Setting.Security.Cors.MaxAgeHour) * time.Hour,
	}))
	r.LoadHTMLGlob(t.Setting.HTTP.TemplatePath)
	r.RedirectTrailingSlash = false
	r.Use(LoggerWithFormatter())
	r.Use(gin.Recovery())
	r.Static(t.Setting.HTTP.MediaURL, t.Setting.HTTP.MediaPath)
	r.Static(t.Setting.HTTP.StaticURL, t.Setting.HTTP.StaticPath)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", func(c *gin.Context) {
		err := Db.DB().Ping()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"Status":     400,
				"APIStatus":  "Error",
				"DBStatus":   "Error",
				"DBError":    err.Error(),
				"APIVersion": t.Setting.Version,
			})
		} else {
			c.JSON(200, gin.H{
				"Status":     200,
				"APIStatus":  "alive",
				"DBStatus":   "alive",
				"DBInfo":     Db.DB().Stats(),
				"APIVersion": t.Setting.Version,
			})
		}
	})
	t.Router = r
}

// NewAPIGroup register new apigroup
func (t *Trinity) NewAPIGroup(path string) *gin.RouterGroup {
	return t.Router.Group(path)
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
		rg.GET("/"+resource+"/:key", viewset)
	}
	if SupportMethodMap["List"] {
		rg.GET("/"+resource, viewset)
	}
	if SupportMethodMap["Create"] {
		rg.POST("/"+resource, viewset)
	}
	if SupportMethodMap["Update"] {
		rg.PATCH("/"+resource+"/:key", viewset)
		rg.PUT("/"+resource+"/:key", viewset)
	}
	if SupportMethodMap["Delete"] {
		rg.DELETE("/"+resource+"/:key", viewset)
	}
}
