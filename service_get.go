package trinity

import (
	"math"
	"strconv"
)

// GetResourceList : List method
func GetResourceList(r *GetMixin) {
	// if r.ViewSetRunTime.BeforeGet != nil {

	// }
	// if r.ViewSetRunTime.Get != nil {

	// }
	// if r.ViewSetRunTime.AfterGet != nil {

	// }

	//Pagination Configure
	//if ispagi true :  return count , next page ,data
	//if ispagi false :  return data only
	var count uint
	PaginationOn := r.ViewSetRunTime.Gcontext.DefaultQuery("PaginationOn", "true")
	switch PaginationOn {
	case "true":
		if err := r.ViewSetRunTime.Db.Scopes(
			r.ViewSetRunTime.DBFilterBackend,
			FilterByFilter(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.FilterByList, r.ViewSetRunTime.FilterCustomizeFunc),
			FilterBySearch(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.SearchingByList),
			QueryBySelect(r.ViewSetRunTime.Gcontext),
			QueryByOrdering(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.OrderingByList),
			QueryByPagination(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.PageSize),
			QueryByPreload(r.ViewSetRunTime.PreloadList),
		).Table(r.ViewSetRunTime.ResourceTableName).Find(r.ViewSetRunTime.ModelSerializerlist).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
			r.ViewSetRunTime.HandleResponse(400, nil, err, ErrLoadDataFailed)
			return
		}
		currentPageNum := r.ViewSetRunTime.Gcontext.DefaultQuery("PageNum", "1")
		currentPageNumInt, err := strconv.Atoi(currentPageNum)
		if err != nil || currentPageNumInt < 0 {
			currentPageNumInt = 1
		}

		PageSizeField := r.ViewSetRunTime.Gcontext.DefaultQuery("PageSize", string(r.ViewSetRunTime.PageSize))
		PageSizeFieldInt, err := strconv.Atoi(PageSizeField)
		if err != nil {
			PageSizeFieldInt = r.ViewSetRunTime.PageSize
		}
		//solve datalist return length =0 and return null problem
		var res map[string]interface{}
		if count > 0 {
			res = map[string]interface{}{
				"data":        r.ViewSetRunTime.ModelSerializerlist,
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
		r.ViewSetRunTime.HandleResponse(200, res, nil, nil)
		return
	default:
		if err := r.ViewSetRunTime.Db.Scopes(
			r.ViewSetRunTime.DBFilterBackend,
			FilterByFilter(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.FilterByList, r.ViewSetRunTime.FilterCustomizeFunc),
			FilterBySearch(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.SearchingByList),
			QueryBySelect(r.ViewSetRunTime.Gcontext),
			QueryByOrdering(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.OrderingByList),
			QueryByPreload(r.ViewSetRunTime.PreloadList),
		).Table(r.ViewSetRunTime.ResourceTableName).Find(r.ViewSetRunTime.ModelSerializerlist).Error; err != nil {
			r.ViewSetRunTime.HandleResponse(400, nil, err, ErrLoadDataFailed)
			return
		}
		r.ViewSetRunTime.HandleResponse(200, r.ViewSetRunTime.ModelSerializerlist, nil, nil)
		return
	}

}
