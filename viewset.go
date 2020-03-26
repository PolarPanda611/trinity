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
			"RETRIEVE": t.setting.GetAtomicRequest(),
			"GET":      t.setting.GetAtomicRequest(),
			"POST":     t.setting.GetAtomicRequest(),
			"PATCH":    t.setting.GetAtomicRequest(),
			"PUT":      t.setting.GetAtomicRequest(),
			"DELETE":   t.setting.GetAtomicRequest(),
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
		PageSize:             t.setting.GetPageSize(),
		EnableOrderBy:        true,
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
		EnableOrderBy:            v.EnableOrderBy,
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

func (v *ViewSetRunTime) loadhttpMethod(c *gin.Context) {
	v.Method = GetRequestType(c)
}

func (v *ViewSetRunTime) loadDB() {
	v.IsAtomicRequest = v.Cfg.AtomicRequestMap[v.Method]
	if v.IsAtomicRequest {
		v.Db = v.Cfg.Db.Begin()
	} else {
		v.Db = v.Cfg.Db
	}
}
func (v *ViewSetRunTime) loadAccessConfig() {
	resourceName := GetTypeName(v.ResourceModel, true)
	if len(v.Cfg.AccessBackendRequireMap) == 0 {
		v.Cfg.AccessBackendRequireMap = map[string][]string{
			"RETRIEVE": []string{"system.view." + resourceName},
			"GET":      []string{"system.view." + resourceName},
			"POST":     []string{"system.add." + resourceName},
			"PATCH":    []string{"system.edit." + resourceName},
			"PUT":      []string{"system.edit." + resourceName},
			"DELETE":   []string{"system.delete." + resourceName},
		}
	}
	v.AccessBackendRequire = v.Cfg.AccessBackendRequireMap[v.Method]
	v.AccessBackendCheck = v.Cfg.AccessBackendCheckMap[v.Method]
}

func (v *ViewSetRunTime) loadLogger() {
	v.DBLogger = &defaultViewRuntimeLogger{ViewRuntime: v}
	v.Db.SetLogger(v.DBLogger)
}

func (v *ViewSetRunTime) loadQuery() {
	v.DBFilterBackend = HandleFilterBackend(v.Cfg, v.Method, v.Gcontext)
	v.FilterBackend = v.Cfg.FilterBackendMap[v.Method]
	v.FilterByList = v.Cfg.FilterByList
	v.FilterCustomizeFunc = v.Cfg.FilterCustomizeFunc
	v.OrderingByList = v.Cfg.OrderingByList
	v.PageSize = v.Cfg.PageSize
	v.PreloadList = v.Cfg.PreloadListMap[v.Method]
	v.SearchingByList = v.Cfg.SearchingByList
	v.EnableOrderBy = v.Cfg.EnableOrderBy
}

func (v *ViewSetRunTime) loadAuthBackend() {
	v.HasAuthCtl = v.Cfg.HasAuthCtl
	v.AuthenticationBackend = v.Cfg.AuthenticationBackendMap[v.Method]
	v.GetCurrentUserAuth = v.Cfg.GetCurrentUserAuth
}

func (v *ViewSetRunTime) loadCallback() {
	v.BeforeRetrieve = v.Cfg.BeforeRetrieve
	v.Retrieve = v.Cfg.Retrieve
	v.AfterRetrieve = v.Cfg.AfterRetrieve

	v.BeforeGet = v.Cfg.BeforeGet
	v.Get = v.Cfg.Get
	v.AfterGet = v.Cfg.AfterGet

	v.PostValidation = v.Cfg.PostValidation
	v.BeforePost = v.Cfg.BeforePost
	v.Post = v.Cfg.Post
	v.AfterPost = v.Cfg.AfterPost

	v.BeforePut = v.Cfg.BeforePut
	v.PutValidation = v.Cfg.PutValidation
	v.Put = v.Cfg.Put
	v.AfterPut = v.Cfg.AfterPut

	v.BeforePatch = v.Cfg.BeforePatch
	v.PatchValidation = v.Cfg.PatchValidation
	v.Patch = v.Cfg.Patch
	v.AfterPatch = v.Cfg.AfterPatch
	v.BeforeDelete = v.Cfg.BeforeDelete
	v.Delete = v.Cfg.Delete
	v.AfterDelete = v.Cfg.AfterDelete
}

func (v *ViewSetRunTime) serviceHandle() {
	// handle request start
	var h ServiceHandler
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

// NewRunTime : new the run time with the default config
// ResourceModel : the main resource model . decide the access authorization and the table name for the resource
// ModelSerializer : the serializer model  , used for retrieve , post and patch service,
// ModelSerializerlist : the serializer model , used for get
func (v *ViewSetCfg) NewRunTime(c *gin.Context, ResourceModel interface{}, ModelSerializer interface{}, ModelSerializerlist interface{}) *ViewSetRunTime {
	vRun := &ViewSetRunTime{
		Gcontext:            c,
		TraceID:             c.GetString("TraceID"),
		ResourceModel:       ResourceModel,
		ResourceTableName:   v.Db.NewScope(ResourceModel).TableName(),
		ModelSerializer:     ModelSerializer,
		ModelSerializerlist: ModelSerializerlist,

		EnableChangeLog:      v.EnableChangeLog,
		EnableDataVersion:    v.EnableDataVersion,
		EnableVersionControl: v.EnableVersionControl,

		Cfg: v,
	}

	vRun.loadhttpMethod(c)
	vRun.loadAuthBackend()
	vRun.loadDB()
	vRun.loadAccessConfig()
	vRun.loadQuery()
	vRun.loadLogger()
	vRun.loadCallback()

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
		if v.GetCurrentUserAuth == nil {
			v.HandleResponse(401, nil, ErrGetUserAuth, ErrGetUserAuth)
			return
		}
		err = v.GetCurrentUserAuth(v.Gcontext, v.Db)
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

	v.serviceHandle()
}
