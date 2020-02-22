package trinity

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

// GetResourceByid : Retrieve method
func GetResourceByid(r *RetrieveMixin) {
	id, err := strconv.ParseInt(r.ViewSetRunTime.Gcontext.Params.ByName("id"), 10, 64)
	if err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, ErrLoadDataFailed)
		return
	}
	if err := r.ViewSetRunTime.Db.Scopes(
		r.ViewSetRunTime.DBFilterBackend,
		FilterByParam(id),
		FilterByFilter(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.FilterByList, r.ViewSetRunTime.FilterCustomizeFunc),
		QueryBySelect(r.ViewSetRunTime.Gcontext),
		QueryByPreload(r.ViewSetRunTime.PreloadList),
	).Table(r.ViewSetRunTime.ResourceTableName).First(r.ViewSetRunTime.ModelSerializer).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			r.ViewSetRunTime.HandleResponse(400, nil, err, gorm.ErrRecordNotFound)
			return

		}
		r.ViewSetRunTime.HandleResponse(400, nil, err, ErrLoadDataFailed)
		return
	}
	r.ViewSetRunTime.HandleResponse(200, r.ViewSetRunTime.ModelSerializer, nil, nil)
	return
}
