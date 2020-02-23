package trinity

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var spiltValue = "__"
var filterCondition = []string{"like", "ilike", "in", "notin", "start", "end", "lt", "lte", "gt", "gte", "isnull", "isempty"}

// FilterQuery for filter query handling
type FilterQuery struct {
	QueryName   string
	QueryValue  string
	tablePrefix string

	// processing
	assosiationParam []string
	queryParam       string
	condition        string
	value            string

	//output
	ConditionSQL string
	ValueSQL     interface{}
}

func (f *FilterQuery) decode() {
	params := strings.Split(f.QueryName, spiltValue)
	paramsLen := len(params)
	if paramsLen == 1 {
		f.queryParam = params[paramsLen-1]
		return
	}
	if paramsLen >= 2 {
		if stringInSlice(params[paramsLen-1], filterCondition) {
			f.condition = params[paramsLen-1]
			f.queryParam = params[paramsLen-2]
			if paramsLen >= 3 {
				f.assosiationParam = params[:paramsLen-2]
			}
		} else {
			f.queryParam = params[paramsLen-1]
			f.assosiationParam = params[:paramsLen-1]
		}
	}
}

// GetFilterConditionSQL get query sql
func (f *FilterQuery) getFilterConditionSQL() {
	switch f.condition {
	case "like":
		f.ConditionSQL = fmt.Sprintf(" %v like ? ", f.queryParam)
		f.ValueSQL = fmt.Sprintf("%v%v%v", "%", f.QueryValue, "%")
		break
	case "ilike":
		f.ConditionSQL = fmt.Sprintf(" %v ilike ? ", f.queryParam)
		f.ValueSQL = fmt.Sprintf("%v%v%v", "%", f.QueryValue, "%")
		break
	case "in":
		f.ConditionSQL = fmt.Sprintf(" %v in (?)  ", f.queryParam)
		f.ValueSQL = strings.Split(f.QueryValue, ",")
		break
	case "notin":
		f.ConditionSQL = fmt.Sprintf(" %v not in (?)  ", f.queryParam)
		f.ValueSQL = strings.Split(f.QueryValue, ",")
		break
	case "start":
		f.ConditionSQL = fmt.Sprintf(" %v  >= ? ", f.queryParam)
		f.ValueSQL = fmt.Sprintf("%v%v", f.QueryValue, " 00:00:00")
		break
	case "end":
		f.ConditionSQL = fmt.Sprintf(" %v  <= ? ", f.queryParam)
		f.ValueSQL = fmt.Sprintf("%v%v", f.QueryValue, " 23:59:59")
		break
	case "isnull":
		f.ConditionSQL = fmt.Sprintf(" %v is not null ", f.queryParam)
		f.ValueSQL = nil
		if f.QueryValue == "true" {
			f.ConditionSQL = fmt.Sprintf(" %v is null ", f.queryParam)
		}
		break
	case "lt":
		f.ConditionSQL = fmt.Sprintf(" %v  < ? ", f.queryParam)
		f.ValueSQL = f.QueryValue
		break
	case "lte":
		f.ConditionSQL = fmt.Sprintf(" %v  <= ? ", f.queryParam)
		f.ValueSQL = f.QueryValue
		break
	case "gt":
		f.ConditionSQL = fmt.Sprintf(" %v  > ? ", f.queryParam)
		f.ValueSQL = f.QueryValue
		break
	case "gte":
		f.ConditionSQL = fmt.Sprintf(" %v  >= ? ", f.queryParam)
		f.ValueSQL = f.QueryValue
		break
	case "isempty":
		f.ConditionSQL = fmt.Sprintf(" (COALESCE(\"%v\"::varchar ,'') != '' )  ", f.queryParam)
		f.ValueSQL = nil
		if f.QueryValue == "true" {
			f.ConditionSQL = fmt.Sprintf(" (COALESCE(\"%v\"::varchar ,'') = '' )  ", f.queryParam)
			f.ValueSQL = nil
		}
		break
	default:
		f.ConditionSQL = fmt.Sprintf(" %v = ? ", f.queryParam)
		f.ValueSQL = f.QueryValue

	}

	// 	break
	return
}

// GetFilterQuerySQL get query
func (f *FilterQuery) GetFilterQuerySQL() {
	f.decode()
	f.getFilterConditionSQL()
	assosiationParamLen := len(f.assosiationParam)
	for range f.assosiationParam {
		lastIndex := assosiationParamLen - 1
		lastParam := f.assosiationParam[lastIndex]
		f.ConditionSQL = fmt.Sprintf(" %v_id in ( select id from %v%v where %v ) ", lastParam, f.tablePrefix, lastParam, f.ConditionSQL)
		assosiationParamLen = assosiationParamLen - 1

	}
	f.ConditionSQL = DeleteExtraSpace(f.ConditionSQL)
	return
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

//FilterByCustomizeCondition filter by customize condition
func FilterByCustomizeCondition(ok bool, k string, v interface{}) func(db *gorm.DB) *gorm.DB {
	if ok {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(k, v)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}

}

//FilterByParam to filter by param return db
func FilterByParam(param int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("\"id\" = ?", param)
	}
}

// HandleFilterBackend handle filter backend
func HandleFilterBackend(v *ViewSetCfg, method string, c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// You can use reqUserID here to check user authorization
		return v.FilterBackendMap[method](c, db)
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
			filter := &FilterQuery{
				QueryName:   queryName,
				QueryValue:  queryValue,
				tablePrefix: GlobalTrinity.setting.Database.TablePrefix,
			}
			filter.GetFilterQuerySQL()
			if filter.ValueSQL == nil {
				db = db.Where(filter.ConditionSQL)
			} else {
				db = db.Where(filter.ConditionSQL, filter.ValueSQL)
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
