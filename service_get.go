package trinity

import (
	"math"
	"strconv"
)

// GetHandler : List method
func GetHandler(r *GetMixin) {
	// if Callback BeforeGet registered , run before get callback
	if IsFuncInited(r.ViewSetRunTime.BeforeGet) {
		r.ViewSetRunTime.BeforeGet(r.ViewSetRunTime)
	}
	// if Callback Get registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.Get) {
		r.ViewSetRunTime.Get(r.ViewSetRunTime)
	}
	// if Callback AfterGet registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.AfterGet) {
		r.ViewSetRunTime.AfterGet(r.ViewSetRunTime)
	}
	if r.ViewSetRunTime.RealError == nil {
		r.ViewSetRunTime.Response()
	}
}

// DefaultGetCallback Default GetHandler
func DefaultGetCallback(r *ViewSetRunTime) {
	//Pagination Configure
	//if ispagi true :  return count , next page ,data
	//if ispagi false :  return data only
	var count uint
	PaginationOn := r.Gcontext.DefaultQuery("PaginationOn", "true")
	switch PaginationOn {
	case "true":
		if err := r.Db.Scopes(
			r.DBFilterBackend,
			FilterByFilter(r.Gcontext, r.FilterByList, r.FilterCustomizeFunc),
			FilterBySearch(r.Gcontext, r.SearchingByList),
			QueryBySelect(r.Gcontext),
			QueryByOrdering(r.Gcontext, r.EnableOrderBy, r.OrderingByList),
			QueryByPagination(r.Gcontext, r.PageSize),
			QueryByPreload(r.PreloadList),
		).Table(r.ResourceTableName).Find(r.ModelSerializerlist).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
			r.HandleResponse(400, nil, err, ErrLoadDataFailed)
			return
		}
		currentPageNum := r.Gcontext.DefaultQuery("PageNum", "1")
		currentPageNumInt, err := strconv.Atoi(currentPageNum)
		if err != nil || currentPageNumInt < 0 {
			currentPageNumInt = 1
		}

		PageSizeField := r.Gcontext.DefaultQuery("PageSize", string(r.PageSize))
		PageSizeFieldInt, err := strconv.Atoi(PageSizeField)
		if err != nil {
			PageSizeFieldInt = r.PageSize
		}
		//solve datalist return length =0 and return null problem
		var res map[string]interface{}
		if count > 0 {
			res = map[string]interface{}{
				"data":        r.ModelSerializerlist,
				"currentpage": currentPageNumInt,
				"totalcount":  count,
				"totalpage":   math.Ceil(float64(count) / float64(PageSizeFieldInt)),
			}
		} else {
			res = map[string]interface{}{
				"data":        []string{},
				"currentpage": currentPageNumInt,
				"totalcount":  count,
				"totalpage":   math.Ceil(float64(count) / float64(PageSizeFieldInt)),
			}
		}
		r.HandleResponse(200, res, nil, nil)
		return
	default:
		if err := r.Db.Scopes(
			r.DBFilterBackend,
			FilterByFilter(r.Gcontext, r.FilterByList, r.FilterCustomizeFunc),
			FilterBySearch(r.Gcontext, r.SearchingByList),
			QueryBySelect(r.Gcontext),
			QueryByOrdering(r.Gcontext, r.EnableOrderBy, r.OrderingByList),
			QueryByPreload(r.PreloadList),
		).Table(r.ResourceTableName).Find(r.ModelSerializerlist).Error; err != nil {
			r.HandleResponse(400, nil, err, ErrLoadDataFailed)
			return
		}
		r.HandleResponse(200, r.ModelSerializerlist, nil, nil)
		return
	}
}
