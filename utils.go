package trinity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CheckFileIsExist : check file if exist ,exist -> true , not exist -> false  ,
/**
 * @param filename string ,the file name need to check
 * @return boolean string
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// GetLogFilePath Initial Log File
func GetLogFilePath(rootpath string, fileName string) string {
	return filepath.Join(rootpath, fileName)

}

// GetRequestType to get http request type with restful style
func GetRequestType(c *gin.Context) string {
	if c.Request.Method == "GET" {
		if len(c.Params) > 0 {
			return "RETRIEVE"
		}
	}
	return c.Request.Method
}

// CheckAccessAuthorization to check access authorization
func CheckAccessAuthorization(requiredPermission, userPermission []string) error {
	if SliceInSlice(requiredPermission, userPermission) {
		return nil
	}
	return ErrAccessAuthCheckFailed
}

// HandleServices for multi response
func HandleServices(m ReqMixinHandler, v *ViewSetRunTime, cHandler ...func(r *ViewSetRunTime)) {
	if len(cHandler) == 1 && cHandler[0] != nil {
		cHandler[0](v)
		return
	}
	m.Handler()
	return
}

//FilterPidByParam filter pid by param
func FilterPidByParam(param string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("\"pid\" = ?", param)
	}
}

//FilterKeyByParam filter key by param
func FilterKeyByParam(param int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("\"id\" = ?", param)
	}
}

//DataVersionFilter filter key by param
func DataVersionFilter(param interface{}, isCheck bool) func(db *gorm.DB) *gorm.DB {
	if isCheck {
		dVersion, _ := param.(string)
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("\"d_version\" = ?", dVersion)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

// HandleFilterBackend handle filter backend
func HandleFilterBackend(v *ViewSetCfg, method string, c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// You can use reqUserID here to check user authorization
		return v.FilterBackendMap[method](c, db)
	}

}

//FilterByParam to filter by param return db
func FilterByParam(param int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("\"id\" = ?", param)
	}
}

//FilterByFilter handle filter
func FilterByFilter(c *gin.Context, FilterByList []string, FilterCustomizeFunc map[string]func(db *gorm.DB, queryValue string) *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Filter condition
		for _, queryName := range FilterByList {
			queryValue := c.Query(queryName)
			if len(queryValue) == 0 {
				continue
			}
			if _, ok := FilterCustomizeFunc[queryName]; ok {
				db = FilterCustomizeFunc[queryName](db, queryValue)
				continue
			}
			if len(strings.Split(queryName, "__")) > 1 {
				switch strings.Split(queryName, "__")[1] {
				case "like":
					db = db.Where("\""+strings.Split(queryName, "__")[0]+"\" like ? ", "%"+queryValue+"%")
					break
				case "ilike":
					db = db.Where("\""+strings.Split(queryName, "__")[0]+"\" ilike ? ", "%"+queryValue+"%")
					break
				case "in":
					strings.Split(queryValue, ",")
					db = db.Where("\""+strings.Split(queryName, "__")[0]+"\" in  (?)", strings.Split(queryValue, ","))
					break
				case "start":
					db = db.Where("\""+strings.Split(queryName, "__")[0]+"\" > ? ", queryValue+" 00:00:00")
					break
				case "end":
					db = db.Where("\""+strings.Split(queryName, "__")[0]+"\" < ? ", queryValue+" 23:59:59")
					break
				case "isnull":
					if queryValue == "true" {
						db = db.Where("( \"" + strings.Split(queryName, "__")[0] + "\" is null )  ")
					}
					if queryValue == "false" {
						db = db.Where("( \"" + strings.Split(queryName, "__")[0] + "\" is not null) ")
					}
					break
				case "isempty":
					if queryValue == "true" {
						db = db.Where("(COALESCE(\"" + strings.Split(queryName, "__")[0] + "\"::varchar ,'') ='' )  ")
					}
					if queryValue == "false" {
						db = db.Where("(COALESCE(\"" + strings.Split(queryName, "__")[0] + "\"::varchar ,'') !='') ")
					}
					break
				case "hasor":
					queryornamelist := strings.Split(strings.Split(queryName, "__")[0], "hasor")
					for i, v := range queryornamelist {
						if i == 0 {
							db = db.Where("( \""+v+"\" = ? ", queryValue)
							continue
						}
						if i == len(queryornamelist)-1 {
							db = db.Or("\""+v+"\" = ? )", queryValue)
							continue
						}
						db = db.Or("\""+v+"\" = ? ", queryValue)

					}
					break
				}
			} else {
				db = db.Where(queryName+" = ? ", queryValue)
			}
		}

		return db
	}
}

//FilterBySearch handle search
func FilterBySearch(c *gin.Context, SearchingByList []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//Searching Section : keywords : Searchby
		SearchValue := c.Query("Searchby")
		if len(SearchValue) != 0 {
			for i, searchField := range SearchingByList {
				if i == 0 {
					db = db.Where("\""+searchField+"\""+" ilike ? ", "%"+SearchValue+"%")
					continue
				}
				db = db.Or("\""+searchField+"\""+" ilike ? ", "%"+SearchValue+"%")

			}
		}
		return db
	}
}

//QueryByOrdering handle ordering
func QueryByOrdering(c *gin.Context, OrderingByList map[string]bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//Searching Section : keywords : OrderingBy , separate by comma
		//example : OrderingBy=-id,user,-name , - means desc , default asc
		OrderByField := c.Query("OrderingBy")
		ordercondition := ""
		for _, orderField := range strings.Split(OrderByField, ",") {
			if len(strings.Split(orderField, "-")) > 1 {
				if _, ok := OrderingByList[strings.Split(orderField, "-")[1]]; ok {
					ordercondition += strings.Split(orderField, "-")[1] + " desc ,"
				}
			} else {
				if _, ok := OrderingByList[orderField]; ok {
					ordercondition += orderField + " asc ,"
				}
			}
		}
		return db.Order(strings.TrimSuffix(ordercondition, ","))
	}
}

//QueryByPagination handling pagination
func QueryByPagination(c *gin.Context, PageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		PageNumField := c.DefaultQuery("PageNum", "0")
		PageSizeField := c.DefaultQuery("PageSize", string(PageSize))
		var err error
		var PageNumFieldInt int
		var PageSizeFieldInt int
		if PageNumFieldInt, err = strconv.Atoi(PageNumField); err != nil || PageNumFieldInt < 0 {
			PageNumFieldInt = 0
		}
		PageNumFieldInt = PageNumFieldInt - 1
		if PageSizeFieldInt, err = strconv.Atoi(PageSizeField); err != nil {
			PageSizeFieldInt = PageSize
		}
		offset := PageNumFieldInt * PageSizeFieldInt
		limit := PageSizeFieldInt
		return db.Offset(offset).Limit(limit)

	}
}

//QueryByPreload handling preload
func QueryByPreload(PreloadList map[string]func(db *gorm.DB) *gorm.DB) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		if len(PreloadList) > 0 {
			for k, v := range PreloadList {
				if v == nil {
					db = db.Preload(k)
				} else {
					db = db.Preload(k, v)
				}

			}
		}
		return db
	}
}

//QueryBySelect handling select
func QueryBySelect(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		SelectByField := c.Query("SelectBy")
		if SelectByField == "" {
			return db
		}
		strings.Split(SelectByField, ",")
		return db.Select(strings.Split(SelectByField, ","))

	}
}

//HandlestructBySelect handling select
func HandlestructBySelect(c *gin.Context, resource interface{}) (int, interface{}, error) {
	SelectByField := c.Query("SelectBy")
	if SelectByField == "" {
		return 200, resource, nil
	}
	m, _ := json.Marshal(&resource)
	// decode it back to get a map
	var a interface{}
	json.Unmarshal(m, &a)
	b := a.(map[string]interface{})
	for k := range b {
		if !stringInSlice(k, strings.Split(SelectByField, ",")) {
			delete(b, k)
		}
	}
	json.Marshal(b)
	return 200, b, nil

}

func handleMultistructBySelect(c *gin.Context, resource interface{}) (int, interface{}, error) {
	SelectByField := c.Query("SelectBy")
	if SelectByField == "" {
		return 200, resource, nil
	}
	return 200, resource, nil

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false

}

//InSlice if value in stringlist
func InSlice(value string, stringSlice []string) bool {
	for _, v := range stringSlice {
		if v == value {
			return true
		}
	}
	return false
}

//SliceInSlice if slice in slice
func SliceInSlice(sliceToCheck []string, slice []string) bool {
	for _, v := range sliceToCheck {
		if !InSlice(v, slice) {
			return false
		}
	}
	return true
}

// GetTypeName to get struct type name
func GetTypeName(myvar interface{}, isToLowerCase bool) string {
	name := ""
	t := reflect.TypeOf(myvar)
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}
	if isToLowerCase {
		name = strings.ToLower(name)
	}
	return name

}

//RecordErrorLevelTwo login error and print line , func , and error to gin context
func RecordErrorLevelTwo() (uintptr, string, int) {
	funcName, file, line, _ := runtime.Caller(2)
	return funcName, file, line
}

// Getparentdirectory : get parent directory of the path ,
/*
 * @param path string  ,the path you want to get parent directory
 * @return string  , the parent directory you need
 */
func Getparentdirectory(path string, level int) string {
	return strings.Join(strings.Split(path, "/")[0:len(strings.Split(path, "/"))-level], "/")
}
