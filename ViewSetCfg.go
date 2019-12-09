package trinity

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ViewSetRunTime : put runtime data
type ViewSetRunTime struct {
	mu     sync.Mutex
	Method string
	// gin.context
	Gcontext *gin.Context
	// db instance
	Db *gorm.DB
	// ResourceModel
	// target resource model
	ResourceModel interface{} // ResourceModel for the resource
	// resource serializer , used to limit the retrieve object
	ModelSerializer interface{}
	// ModelSerializerlist
	// resource serializer , used to limit the get object list
	ModelSerializerlist interface{}
	// HasAuthCtl
	// if do the auth check ,default false
	HasAuthCtl            bool
	AuthenticationBackend func(c *gin.Context) error
	GetCurrentUserAuth    interface{}
	AccessBackendRequire  []string
	AccessBackendCheck    func(v *ViewSetRunTime) error
	DBFilterBackend       func(db *gorm.DB) *gorm.DB // current dbfilterbackend
	PreloadList           []string
	FilterBackend         func(c *gin.Context, db *gorm.DB) *gorm.DB
	FilterByList          []string
	FilterCustomizeFunc   map[string]func(db *gorm.DB, queryValue string) *gorm.DB
	SearchingByList       []string
	OrderingByList        map[string]bool
	PageSize              int
	Retrieve              func(r *ViewSetRunTime)
	Get                   func(r *ViewSetRunTime)
	Post                  func(r *ViewSetRunTime)
	Put                   func(r *ViewSetRunTime)
	Patch                 func(r *ViewSetRunTime)
	Delete                func(r *ViewSetRunTime)
	Cfg                   *ViewSetCfg

	//response handle
	Status    int
	ResBody   interface{}
	FuncName  uintptr
	File      string
	Line      int
	RealError error
	UserError error
}

// ViewSetCfg for viewset config
type ViewSetCfg struct {
	sync.RWMutex
	// global config
	Db *gorm.DB
	// HasAuthCtl
	// if do the auth check ,default false
	HasAuthCtl bool
	// AuthenticationBackendMap
	// if HasAuthCtl == false ; pass... customize the authentication check , default jwt  ;
	// please set UserID in context
	// e.g : c.Set("UserID", tokenClaims.UID)
	AuthenticationBackendMap map[string]func(c *gin.Context) error
	// GetCurrentUserAuth
	// must be type : func(c *gin.Context, db *gorm.DB) error
	// if HasAuthCtl == false ; pass...
	// get user auth func with UserID if you set in AuthenticationBackend
	// please set UserPermission and UserKey in context
	// e.g : c.Set("UserKey",UserKey) with c.GetString("UserID")
	// e.g : c.Set("UserPermission", UserPermission) with c.GetString("UserID")
	GetCurrentUserAuth interface{}
	// AccessBackendReqMap
	// if HasAuthCtl == false ; pass... customize the access require permission
	AccessBackendRequireMap map[string][]string
	// AccessBackendCheckMap
	// if HasAuthCtl == false ; pass... customize the access check , check user permission
	// e.g : userPermission :=  c.GetString("UserPermission")
	// e.g : requiredPermission := []string{"123"} get with AccessBackendReqMap by default
	// e.g : trinity.CheckAccessAuthorization(requiredPermission , userPermission) , true?allow:deny
	AccessBackendCheckMap map[string]func(v *ViewSetRunTime) error
	// PreloadListMap gorm preload list
	PreloadListMap map[string][]string
	// FilterBackendMap : all the query will with this filter backend
	FilterBackendMap map[string]func(c *gin.Context, db *gorm.DB) *gorm.DB
	// FilterByList : only in FilterByList will do the filter
	FilterByList []string
	// FilterCustomizeFunc : can do the customize filter ,mapping with FilterByList
	FilterCustomizeFunc map[string]func(db *gorm.DB, queryValue string) *gorm.DB
	// SearchingByList : with keyword "SearchBy" on url query ,
	// will do the where (xxx =? or xxx=?)
	SearchingByList []string
	// OrderingByList : with keyword "OrderBy" on url query ,
	// only define in OrderingByList will do the order by
	// e.g: OrderBy=xxx-   ==> order by xxx desc
	// e.g: OrderBy=xxx   ==> order by xxx asc
	OrderingByList map[string]bool
	// PageSize default 10
	// keyword : PageNum , PageSize to do the limit and offset
	PageSize int
	// Retrieve: customize retrieve func
	Retrieve func(r *ViewSetRunTime)
	// Get: customize Get func
	Get func(r *ViewSetRunTime)
	// Post: customize Post func
	Post func(r *ViewSetRunTime)
	// Put: customize Put func
	Put func(r *ViewSetRunTime)
	// Patch: customize Patch func
	Patch func(r *ViewSetRunTime)
	// Delete: customize Delete func
	Delete func(r *ViewSetRunTime)
}
