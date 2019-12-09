package trinity

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

// GetResourceByid : Retrieve method
func GetResourceByid(r *RetrieveMixin) {
	if err := r.ViewSetRunTime.Db.Scopes(
		r.ViewSetRunTime.DBFilterBackend,
		FilterByParam(r.ViewSetRunTime.Gcontext.Params.ByName("key")),
		FilterByFilter(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.FilterByList, r.ViewSetRunTime.FilterCustomizeFunc),
		QueryBySelect(r.ViewSetRunTime.Gcontext),
		QueryByPreload(r.ViewSetRunTime.PreloadList),
	).Table(r.ViewSetRunTime.Db.NewScope(r.ResourceModel).TableName()).First(r.ViewSetRunTime.ModelSerializer).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			r.ViewSetRunTime.HandleResponse(400, nil, err, errors.New("app.error.recordnotfound"))
			return

		}
		r.ViewSetRunTime.HandleResponse(400, nil, err, errors.New("app.err.loaddatafailed"))
		return
	}
	r.ViewSetRunTime.HandleResponse(200, r.ModelSerializer, nil, nil)
	return
}

// GetResourceList : List method
func GetResourceList(r *GetMixin) {

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
		).Table(r.ViewSetRunTime.Db.NewScope(r.ResourceModel).TableName()).Find(r.ViewSetRunTime.ModelSerializerlist).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
			r.ViewSetRunTime.HandleResponse(400, nil, err, errors.New("app.err.loaddatalistfailed"))
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
		).Table(r.ViewSetRunTime.Db.NewScope(r.ResourceModel).TableName()).Find(r.ViewSetRunTime.ModelSerializerlist).Error; err != nil {
			r.ViewSetRunTime.HandleResponse(400, nil, err, errors.New("app.err.loaddatalistfailed"))
			return
		}
		r.ViewSetRunTime.HandleResponse(200, r.ViewSetRunTime.ModelSerializerlist, nil, nil)
		return
	}

}

// CreateResource : Create method
func CreateResource(r *PostMixin) {
	if err := r.ViewSetRunTime.Gcontext.BindJSON(r.ViewSetRunTime.ModelSerializer); err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, errors.New("app.err.resolvedataerror"))
		return
	}
	if err := r.ViewSetRunTime.Db.Table(r.ViewSetRunTime.Db.NewScope(r.ResourceModel).TableName()).Create(r.ViewSetRunTime.ModelSerializer).Error; err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, errors.New("app.err.createdatalistfailed"))
		return
	}
	r.ViewSetRunTime.HandleResponse(201, r.ViewSetRunTime.ModelSerializer, nil, nil)
	return
}

// PatchResource : PATCH method
func PatchResource(r *PatchMixin) {
	buf := make([]byte, 1024)
	n, _ := r.ViewSetRunTime.Gcontext.Request.Body.Read(buf)
	requestbodyMap := make(map[string]interface{})
	if err := json.Unmarshal(buf[0:n], &requestbodyMap); err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, errors.New("app.err.resolvedataerror"))
		return
	}
	if err := r.ViewSetRunTime.Db.Scopes(
		r.ViewSetRunTime.DBFilterBackend,
		FilterByParam(r.ViewSetRunTime.Gcontext.Params.ByName("key")),
	).Table(r.ViewSetRunTime.Db.NewScope(r.ResourceModel).TableName()).Updates(requestbodyMap).First(r.ViewSetRunTime.ModelSerializer).Error; err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, errors.New("app.err.updatedatalistfailed"))
		return
	}
	r.ViewSetRunTime.HandleResponse(200, r.ViewSetRunTime.ModelSerializer, nil, nil)
	return

}

// DeleteResource : DELETE method
func DeleteResource(r *DeleteMixin) {
	requestbodyMap := make(map[string]interface{})
	nowTime := time.Now()
	requestbodyMap["deleted_time"] = nowTime
	requestbodyMap["delete_user_key"] = r.ViewSetRunTime.Gcontext.GetString("UserKey")
	if err := r.ViewSetRunTime.Db.Scopes(
		r.ViewSetRunTime.DBFilterBackend,
		FilterByParam(r.ViewSetRunTime.Gcontext.Params.ByName("key")),
	).Table(r.ViewSetRunTime.Db.NewScope(r.ResourceModel).TableName()).Updates(requestbodyMap).Error; err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, errors.New("app.err.deletedatalistfailed"))
		return

	}
	r.ViewSetRunTime.HandleResponse(200, "Delete Successfully", nil, nil)
	return
}
