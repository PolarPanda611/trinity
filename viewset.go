package trinity

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// initViewSetCfg for  initial ViewSetCfg
func (t *Trinity) initViewSetCfg() {
	v := &ViewSetCfg{
		Db:         t.db,
		HasAuthCtl: false,
		AtomicRequestMap: map[string]bool{
			"RETRIEVE": t.setting.Webapp.AtomicRequest,
			"GET":      t.setting.Webapp.AtomicRequest,
			"POST":     t.setting.Webapp.AtomicRequest,
			"PATCH":    t.setting.Webapp.AtomicRequest,
			"PUT":      t.setting.Webapp.AtomicRequest,
			"DELETE":   t.setting.Webapp.AtomicRequest,
		},
		AuthenticationBackendMap: map[string]func(c *gin.Context) error{
			"RETRIEVE": JwtUnverifiedAuthBackend,
			"GET":      JwtUnverifiedAuthBackend,
			"POST":     JwtUnverifiedAuthBackend,
			"PATCH":    JwtUnverifiedAuthBackend,
			"PUT":      JwtUnverifiedAuthBackend,
			"DELETE":   JwtUnverifiedAuthBackend,
		},
		GetCurrentUserAuth: func(c *gin.Context, db *gorm.DB) error {
			c.Set("UserID", "")                 // with  c.GetInt64("UserID")
			c.Set("UserPermission", []string{}) // with  c.GetStringSlice("UserID")
			return nil
		},
		AccessBackendRequireMap: map[string][]string{},
		AccessBackendCheckMap: map[string]func(v *ViewSetRunTime) error{
			"RETRIEVE": DefaultAccessBackend,
			"GET":      DefaultAccessBackend,
			"POST":     DefaultAccessBackend,
			"PATCH":    DefaultAccessBackend,
			"PUT":      DefaultAccessBackend,
			"DELETE":   DefaultAccessBackend,
		},
		PreloadListMap: map[string]map[string]func(db *gorm.DB) *gorm.DB{
			"RETRIEVE": nil, // Foreign key :Foreign table  if you want to filter the foreign table
			"GET":      nil,
		},
		FilterBackendMap: map[string]func(c *gin.Context, db *gorm.DB) *gorm.DB{
			"RETRIEVE": DefaultFilterBackend,
			"GET":      DefaultFilterBackend,
			"POST":     DefaultFilterBackend,
			"PATCH":    DefaultFilterBackend,
			"PUT":      DefaultFilterBackend,
			"DELETE":   DefaultFilterBackend,
		},
		FilterByList:         []string{},
		FilterCustomizeFunc:  map[string]func(db *gorm.DB, queryValue string) *gorm.DB{},
		SearchingByList:      []string{},
		OrderingByList:       map[string]bool{},
		PageSize:             t.setting.Webapp.PageSize,
		EnableChangeLog:      false,
		EnableDataVersion:    true,
		EnableVersionControl: false,
		Retrieve:             DefaultRetrieveCallback,
		Get:                  DefaultGetCallback,
		Post:                 DefaultPostCallback,
		Put:                  DefaultPutCallback,
		Patch:                DefaultPatchCallback,
		Delete:               DefaultDeleteCallback,
	}
	t.vCfg = v

}

func (v *ViewSetCfg) clone() *ViewSetCfg {
	vClone := &ViewSetCfg{
		Db:                       v.Db,
		HasAuthCtl:               v.HasAuthCtl,
		AtomicRequestMap:         v.AtomicRequestMap,
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
		EnableChangeLog:          v.EnableChangeLog,
		EnableDataVersion:        v.EnableDataVersion,
		EnableVersionControl:     v.EnableVersionControl,
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
// ResourceModel : the main resource model . decide the access authorization and the table name for the resource
// ModelSerializer : the serializer model  , used for retrieve , post and patch service,
// ModelSerializerlist : the serializer model , used for get
func (v *ViewSetCfg) NewRunTime(c *gin.Context, ResourceModel interface{}, ModelSerializer interface{}, ModelSerializerlist interface{}) *ViewSetRunTime {
	httpMethod := GetRequestType(c)
	resourceName := GetTypeName(ResourceModel, true)
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
	var db *gorm.DB
	isAtomicRequest := v.AtomicRequestMap[httpMethod]
	if isAtomicRequest {
		db = v.Db.Begin()
	} else {
		db = v.Db
	}

	vRun := &ViewSetRunTime{
		Gcontext:              c,
		TraceID:               c.GetString("TraceID"),
		IsAtomicRequest:       isAtomicRequest,
		Db:                    db,
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
		EnableChangeLog:       v.EnableChangeLog,
		EnableDataVersion:     v.EnableDataVersion,
		EnableVersionControl:  v.EnableVersionControl,

		BeforeRetrieve: v.BeforeRetrieve,
		Retrieve:       v.Retrieve,
		AfterRetrieve:  v.AfterRetrieve,

		BeforeGet: v.BeforeGet,
		Get:       v.Get,
		AfterGet:  v.AfterGet,

		PostValidation: v.PostValidation,
		BeforePost:     v.BeforePost,
		Post:           v.Post,
		AfterPost:      v.AfterPost,

		BeforePut:     v.BeforePut,
		PutValidation: v.PutValidation,
		Put:           v.Put,
		AfterPut:      v.AfterPut,

		BeforePatch:     v.BeforePatch,
		PatchValidation: v.PatchValidation,
		Patch:           v.Patch,
		AfterPatch:      v.AfterPatch,
		BeforeDelete:    v.BeforeDelete,
		Delete:          v.Delete,
		AfterDelete:     v.AfterDelete,
		Cfg:             v,
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

	// first level : authentication control
	if v.HasAuthCtl {
		err := v.AuthenticationBackend(v.Gcontext)
		// if err return 401 unauthorized
		if err != nil {
			v.HandleResponse(401, nil, err, ErrUnverifiedToken)
			return
		}
		// get user auth
		getCurrentUserAuth, ok := v.GetCurrentUserAuth.(func(c *gin.Context, db *gorm.DB) error)
		if !ok {
			v.HandleResponse(401, nil, ErrGetUserAuth, ErrGetUserAuth)
			return
		}
		err = getCurrentUserAuth(v.Gcontext, v.Db)
		if err != nil {
			v.HandleResponse(401, nil, err, ErrGetUserAuth)
			return
		}
		// Access control
		if err := v.AccessBackendCheck(v); err != nil {
			v.HandleResponse(401, nil, err, ErrAccessAuthCheckFailed)
			return
		}
	}
	v.Db = v.Db.Set("UserID", v.Gcontext.GetInt64("UserID"))

	// handle request start
	var h ReqMixinHandler
	switch v.Method {
	case "RETRIEVE":
		h = &RetrieveMixin{v}
		break
	case "GET":
		h = &GetMixin{v}
		break
	case "POST":
		h = &PostMixin{v}
		break
	case "PUT":
		h = &PutMixin{v}
		break
	case "PATCH":
		h = &PatchMixin{v}
		break
	case "DELETE":
		h = &DeleteMixin{v}
		break
	default:
		h = &UnknownMixin{v}
		break
	}
	h.Handler()
	return

}
