package trinity

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// InitViewSetCfg for  initial ViewSetCfg
func (t *Trinity) InitViewSetCfg() {
	v := &ViewSetCfg{
		Db:         t.db,
		HasAuthCtl: false,
		AuthenticationBackendMap: map[string]func(c *gin.Context) (error, error){
			"RETRIEVE": JwtUnverifiedAuthBackend,
			"GET":      JwtUnverifiedAuthBackend,
			"POST":     JwtUnverifiedAuthBackend,
			"PATCH":    JwtUnverifiedAuthBackend,
			"PUT":      JwtUnverifiedAuthBackend,
			"DELETE":   JwtUnverifiedAuthBackend,
		},
		GetCurrentUserAuth: func(c *gin.Context, db *gorm.DB) (error, error) {
			c.Set("UserKey", "")                // with  c.GetString("UserID")
			c.Set("UserPermission", []string{}) // with  c.GetString("UserID")
			return nil, nil
		},
		AccessBackendRequireMap: map[string][]string{},
		AccessBackendCheckMap: map[string]func(v *ViewSetRunTime) (error, error){
			"RETRIEVE": DefaultAccessBackend,
			"GET":      DefaultAccessBackend,
			"POST":     DefaultAccessBackend,
			"PATCH":    DefaultAccessBackend,
			"PUT":      DefaultAccessBackend,
			"DELETE":   DefaultAccessBackend,
		},
		PreloadListMap: map[string][]string{
			"RETRIEVE": []string{}, // Foreign key :Foreign table  if you want to filter the foreign table
			"GET":      []string{},
		},
		FilterBackendMap: map[string]func(c *gin.Context, db *gorm.DB) *gorm.DB{
			"RETRIEVE": DefaultFilterBackend,
			"GET":      DefaultFilterBackend,
			"POST":     DefaultFilterBackend,
			"PATCH":    DefaultFilterBackend,
			"PUT":      DefaultFilterBackend,
			"DELETE":   DefaultFilterBackend,
		},
		FilterByList:        []string{},
		FilterCustomizeFunc: map[string]func(db *gorm.DB, queryValue string) *gorm.DB{},
		SearchingByList:     []string{},
		OrderingByList:      map[string]bool{},
		PageSize:            10,
	}
	t.Lock()
	t.vCfg = v
	t.Unlock()

}

func (v *ViewSetCfg) clone() *ViewSetCfg {
	vClone := &ViewSetCfg{
		Db:                       v.Db,
		HasAuthCtl:               v.HasAuthCtl,
		AuthenticationBackendMap: v.AuthenticationBackendMap,
		GetCurrentUserAuth:       v.GetCurrentUserAuth,
		AccessBackendRequireMap:  v.AccessBackendRequireMap,
		AccessBackendCheckMap:    v.AccessBackendCheckMap,
		PreloadListMap:           v.PreloadListMap,
		FilterBackendMap:         v.FilterBackendMap,
		FilterByList:             v.FilterByList,
		FilterCustomizeFunc:      v.FilterCustomizeFunc,
		SearchingByList:          v.SearchingByList,
		OrderingByList:           v.OrderingByList,
		PageSize:                 v.PageSize,
		Retrieve:                 v.Retrieve,
		Get:                      v.Get,
		Post:                     v.Post,
		Put:                      v.Put,
		Patch:                    v.Patch,
		Delete:                   v.Delete,
	}
	return vClone
}

// NewViewSet new api viewset
func NewViewSet() *ViewSetCfg {
	v := GlobalTrinity.vCfg.clone()
	return v
}

// NewRunTime : new the run time with the default config
func (v *ViewSetCfg) NewRunTime(c *gin.Context, ResourceModel interface{}, ModelSerializer interface{}, ModelSerializerlist interface{}) *ViewSetRunTime {
	httpMethod := GetRequestType(c)
	resourceName := GetTypeName(ResourceModel)
	if len(v.AccessBackendRequireMap) == 0 {
		v.AccessBackendRequireMap = map[string][]string{
			"RETRIEVE": []string{"system.view." + resourceName},
			"GET":      []string{"system.view." + resourceName},
			"POST":     []string{"system.add." + resourceName},
			"PATCH":    []string{"system.edit." + resourceName},
			"PUT":      []string{"system.edit." + resourceName},
			"DELETE":   []string{"system.delete." + resourceName},
		}
	}

	vRun := &ViewSetRunTime{
		Gcontext:              c,
		TraceID:               c.GetString("TraceID"),
		Db:                    v.Db,
		Method:                httpMethod,
		ResourceModel:         ResourceModel,
		ResourceTableName:     v.Db.NewScope(ResourceModel).TableName(),
		ModelSerializer:       ModelSerializer,
		ModelSerializerlist:   ModelSerializerlist,
		HasAuthCtl:            v.HasAuthCtl,
		AuthenticationBackend: v.AuthenticationBackendMap[httpMethod],
		GetCurrentUserAuth:    v.GetCurrentUserAuth,
		AccessBackendRequire:  v.AccessBackendRequireMap[httpMethod],
		AccessBackendCheck:    v.AccessBackendCheckMap[httpMethod],
		DBFilterBackend:       HandleFilterBackend(v, httpMethod, c),
		FilterBackend:         v.FilterBackendMap[httpMethod],
		FilterByList:          v.FilterByList,
		FilterCustomizeFunc:   v.FilterCustomizeFunc,
		OrderingByList:        v.OrderingByList,
		PageSize:              v.PageSize,
		PreloadList:           v.PreloadListMap[httpMethod],
		SearchingByList:       v.SearchingByList,
		Retrieve:              v.Retrieve,
		Get:                   v.Get,
		Post:                  v.Post,
		Put:                   v.Put,
		Patch:                 v.Patch,
		Delete:                v.Delete,
		Cfg:                   v,
	}
	vRun.DBLogger = &defaultViewRuntimeLogger{ViewRuntime: vRun}
	vRun.Db.SetLogger(vRun.DBLogger)
	return vRun
}

// ViewSetServe for viewset handle
// serve flow
// if HasAuthCtl == false {
// 	1.AuthenticationBackend : do the authentication check , normally get the user identity
// 	2.GetCurrentUserAuth : get the user permission information
// 	3.AccessBackend : do the access check
// }
// 4.request data validation
// 5.DbWithBackend : do the DB backend check
// 6.do the request
func (v *ViewSetRunTime) ViewSetServe() {
	//set default value for viewCfg
	// handle request start
	var h ReqMixinHandler
	var ch func(r *ViewSetRunTime) // Customize Handler
	// first level : authentication control
	if v.HasAuthCtl {
		rErr, uErr := v.AuthenticationBackend(v.Gcontext)
		// if err return 401 unauthorized
		if rErr != nil {
			v.HandleResponse(401, nil, rErr, uErr)
			return
		}
		// get user auth
		getCurrentUserAuth, ok := v.GetCurrentUserAuth.(func(c *gin.Context, db *gorm.DB) (error, error))
		if !ok {
			v.HandleResponse(401, nil, ErrGetUserAuth, ErrGetUserAuth)
			return
		}
		rErr, uErr = getCurrentUserAuth(v.Gcontext, v.Db)
		if rErr != nil {
			v.HandleResponse(401, nil, rErr, uErr)
			return
		}
		// Access control
		if rErr, uErr := v.AccessBackendCheck(v); rErr != nil {
			v.HandleResponse(401, nil, rErr, uErr)
			return
		}
	}
	v.Db = v.Db.Set("reqUserKey", v.Gcontext.GetString("reqUserKey"))

	switch v.Method {
	case "RETRIEVE":
		// Retrieve request
		h = &RetrieveMixin{v}
		ch = v.Retrieve
		break
	case "GET":
		// Get request
		h = &GetMixin{v}
		ch = v.Get
		break
	case "POST":
		// Post request
		h = &PostMixin{v}
		ch = v.Post
		break
	case "PUT":
		// Put request
		h = &PutMixin{v}
		ch = v.Put
		break
	case "PATCH":
		// Patch request
		h = &PatchMixin{v}
		ch = v.Patch
		break
	case "DELETE":
		// Delete request
		h = &DeleteMixin{v}
		ch = v.Delete
		break
	default:
		// Unknown request
		h = &UnknownMixin{v}
		break
	}
	HandleServices(h, v, ch)
	return

}
