# trinity

[![Build Status](https://travis-ci.org/PolarPanda611/trinity.svg)](https://travis-ci.org/PolarPanda611/trinity)
[![GitHub Actions](https://github.com/aofei/air/workflows/Main/badge.svg)](https://github.com/PolarPanda611/trinity)
[![codecov](https://codecov.io/gh/aofei/air/branch/master/graph/badge.svg)](https://codecov.io/gh/PolarPanda611/trinity)
[![Go Report Card](https://goreportcard.com/badge/github.com/PolarPanda611/trinity)](https://goreportcard.com/report/github.com/PolarPanda611/trinity)
[![GoDoc](https://godoc.org/github.com/PolarPanda611/trinity?status.svg)](https://godoc.org/github.com/PolarPanda611/trinity)
[![Release](https://img.shields.io/github/release/PolarPanda611/trinity.svg?style=flat-square)](https://github.com/PolarPanda611/trinity/releases)

golang restframework plugin with gin+gorm, 5 lines to generate rest api with high security    
p.s: django restframework like :)

## 特性

* 集成gorm
* 集成gin
* 快速注册路由
* 链路追踪，返回错误代码行数,及sql
* 快速migrate表并生成权限
* 支持快速开发rest风格api并实现增删改查，开箱即用
* JWT token生成，认证，刷新中间件
* json log风格中间件
* 支持请求事务
* 支持自定义用户权限查询
* 支持自定义接口访问权限
* 支持自定义接口数据访问权限
* 自定义分页
* 自定义过滤查询
* 自定义搜索
* 自定义预加载（gorm preload）
* 自定义排序
* 自定义查询包含字段
* 自定义http方法重写
* 添加了change log


## 安装

打开终端输入

```bash
$ go get -u github.com/PolarPanda611/trinity
```

done.

## 文档
```ViewSetCfg.go

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
	Retrieve func(r *ViewSetRunTime) *ViewSetRunTime
	// Get: customize Get func
	Get func(r *ViewSetRunTime) *ViewSetRunTime
	// Post: customize Post func
	Post func(r *ViewSetRunTime) *ViewSetRunTime
	// Put: customize Put func
	Put func(r *ViewSetRunTime) *ViewSetRunTime
	// Patch: customize Patch func
	Patch func(r *ViewSetRunTime) *ViewSetRunTime
	// Delete: customize Delete func
	Delete func(r *ViewSetRunTime) *ViewSetRunTime
}
```


## 例子

* 准备条件,建立model
```model.go
//Country model Country
type Country struct {
	Model
	Code        string `json:"code" gorm:"type:varchar(50);index;unique;not null;"`
	Name        string `json:"name" gorm:"type:varchar(50);"`
	Description string `json:"description" gorm:"type:varchar(100);index;"`
}
```
*初始化trinity设置
```trinityinit/trinityinit.go
// GlobalViewSet global view set config by default value
var GlobalViewSet *trinity.ViewSetCfg


//Inittrinity get default setting
func Inittrinity() {
	trinity.Jwtexpirehour = setting.Cfg.Jwt.Jwtexpirehour
	// Jwtheaderprefix for jwt
	trinity.Jwtheaderprefix = setting.Cfg.Jwt.Jwtheaderprefix
	// Secretkey for jwt
	trinity.Secretkey = setting.Cfg.Secretkey
	// Jwtissuer for jwt
	trinity.Jwtissuer = setting.Cfg.Jwt.Jwtissuer

	GlobalViewSet = trinity.InitDefault(db.Dsource)

	// 在此覆盖默认全局设置
	GlobalViewSet.PageSize = setting.Cfg.Pagesize
}

```

* 使用log，jwt中间件及注册路由

```router.go
func Router() *gin.Engine {
    r := gin.New()
    ````
    r.Use(trinity.JWT())                 // jwt中间件
    r.Use(trinity.LoggerWithFormatter()) // log中间件,链路追踪需开启log中间件，获取TraceID *gin.Context.GetString("TraceID")
    v1 := r.Group("/api/v1")
    {
        // register RESTFUL API router by resouce name
	/*
	*@param RouterGroup :  the router group you want to register
	*@param Resource : the resource of the REST API
	*@param ViewSet : the service of the REST API
	*@param SupportedMethod : the service list of the REST API
	 */
	 // same as 
	 //r.GET("/"+resource+"/:key", viewset)
	 //r.GET("/"+resource, viewset)
	 //r.POST("/"+resource, viewset)
	 //r.PATCH("/"+resource+"/:key", viewset)
	 //r.PUT("/"+resource+"/:key", viewset)
	 //r.DELETE("/"+resource+"/:key", viewset)
        trinity.RegisterRestStyleRouter(v1, "users", servicev1.UserViewSet, []string{"Retrieve", "List", "Create", "Update", "Delete"})
    }
    ```
```

* 注册路由处理
```services/v1/apperror.go
// CountryViewSet hanlde router
func CountryViewSet(c *gin.Context) {
	v := trinityinit.GlobalViewSet.New()
	// 在此覆盖默认局部设置
	v.HasAuthCtl = true
	v.FilterByList = []string{"trace_id"}
	// v.GetCurrentUserAuth = func(c *gin.Context, db *gorm.DB) error {
	// 	c.Set("UserPermission", []string{"system.view.Country"})
	// 	return nil
	// }
	utils.HandleResponse(v.NewRunTime(
		c,
		&model.Country{},
		&model.Country{},
		&[]model.Country{},
	).ViewSetServe())
}
```
* 支持关键字


   * 过滤 (系统预置关键字)
```
	
	//Example :http://127.0.0.1/countries?xxx__like=ooo&xxx__ilike=ooo&xxx__start=date1
	----FilterByList   filter condition must configured in FilterByList config  
	// v.FilterByList=[]string{"xxx__like","xxx__ilike"}
	xxx__like=ooo 		=> xxx like '%ooo%'  
	xxx__ilike=ooo 		=> xxx like '%ooo%'  caps no sensitive
	xxx__in=aaa,bbb,ccc 	=> xxx in ['aaa','bbb','ccc']  
	xxx__start=date1	=> xxx > date1  
	xxx__end=date2		=> xxx < date2
	xxx__isnull=true 	=> xxx is null 
	xxx__isnull=false 	=> xxx is not null 
	xxx__isempty=true	=>(COALESCE("xxx"::varchar ,'') ='' )  
	xxx__isempty=false	=>(COALESCE("xxx"::varchar ,'') !='' )  
	xxxhasoraaa=fff		=> (xxx = fff or aaa = fff)
	xxx=ooo			=>xxx = ooo
```

   * 自定义过滤 (自定义关键字和查询方法)
```
	#Example  :http://127.0.0.1/countries?filterpolarpanda=me
	----FilterByList   filter condition must configured in FilterByList config  
	// v.FilterByList=[]string{"filterpolarpanda"}
	// v.FilterCustomizeFunc=map[string]func(db *gorm.DB, queryValue string) *gorm.DB{
		"filterpolarpanda":func(db *gorm.DB, queryValue string) *gorm.DB{
			db.Where("polarpanda = ?" , "me")
		}
	}
	filterpolarpanda=me 		=> polarpanda = 'me'  
	
```

   * 自定义搜索 (自定义关键字集合)
```
	#Example  :http://127.0.0.1/countries?SearchBy=xxx
	----SearchingByList   search condition must configured in in SearchingByList config  
	//v.SearchingByList=[]string{"name","address"}
	Searchby=xxx		=>(name ilike '%xxx%' or address ilike '%xxx%')  
```

   * 自定义排序 (自定义排序关键字)
```
	#Example  :http://127.0.0.1/countries?OrderingBy=-id,name
	----OrderingByList  order by condition must configured in in OrderingByList config  
	OrderingBy=-id,name	=> order by id desc , name    
```

   * 自定义分页 (自定义分页数量和页码)
```
	//默认开启分页
	#Example  :http://127.0.0.1/countries?PageSize=44&PageNum=2
	----queryByPagination By default , the list will be paged , default page size is configured in Pagination config  
	//offset := PageNumFieldInt * PageSizeFieldInt -1 
	//limit := PageSizeFieldInt  
	PageNum=1		=>PageNum= (1,2.....)  
	PageSize=10		=>PageSize= default value:10
	PaginationOn		=>by default :true , will open the pagination , if else , close the pagination, return all the list 
	//关闭分页
	#Example  :http://127.0.0.1/countries?PaginationOn=false
```
   * 自定义请求回复 
```

// ResponseData http response
type ResponseData struct {
	Status  int         // the http response status  to return
	Result  interface{} // the response data  if req success
	TraceID string
}

// Response handle trinity return value
func Response(r *ViewSetRunTime) {
	var res ResponseData
	res.Status = r.Status
	res.TraceID = r.Gcontext.GetString("TraceID")
	if r.RealError != nil {
		r.Cfg.Logger.LogWriter(r)
		r.Gcontext.Error(r.RealError)
		r.Gcontext.Error(r.UserError)
		res.Result = r.UserError.Error()
		r.Gcontext.AbortWithStatusJSON(r.Status, res)
	} else {
		res.Result = r.ResBody
		r.Gcontext.JSON(r.Status, res)
	}
	return

}
```
   * 支持请求事务
```
	local:
		project: xx 
		version: xxx
		runtime:
			debug: True
		security:
			...
		webapp:
			...
			atomicrequest: true
``` 

   * 添加change log
   目前不支持嵌套map添加change log
```
	migrate  		&trinity.AppChangelog{}, 

	在viewset中
	// PartsViewSet hanlde router
	func PartsViewSet(c *gin.Context) {

		v := trinity.NewViewSet()
		v.HasAuthCtl = true
		v.EnableChangeLog = true
		v.NewRunTime(
			c,
			&model.Part{},
			&model.Part{},
			&[]model.Part{},
		).ViewSetServe()
	}

``` 


## to do list : 
支持字段级验证请求数据   
外键字段过滤   


