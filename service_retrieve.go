package trinity

import (
	"strconv"

	"github.com/PolarPanda611/gorm"
)

// RetriveHandler : List method
func RetriveHandler(r *RetrieveMixin) {
	// if Callback BeforeGet registered , run before get callback
	if IsFuncInited(r.ViewSetRunTime.BeforeRetrieve) {
		r.ViewSetRunTime.BeforeRetrieve(r.ViewSetRunTime)
	}
	// if Callback Get registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.Retrieve) {
		r.ViewSetRunTime.Retrieve(r.ViewSetRunTime)
	}
	// if Callback AfterGet registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.AfterRetrieve) {
		r.ViewSetRunTime.AfterRetrieve(r.ViewSetRunTime)
	}
	if r.ViewSetRunTime.RealError == nil {
		r.ViewSetRunTime.Response()
	}

}

// DefaultRetrieveCallback : Retrieve method
func DefaultRetrieveCallback(r *ViewSetRunTime) {
	id, err := strconv.ParseInt(r.Gcontext.Params.ByName("id"), 10, 64)
	if err != nil {
		r.HandleResponse(400, nil, err, ErrLoadDataFailed)
		return
	}
	if err := r.Db.Scopes(
		r.DBFilterBackend,
		FilterByParam(id),
		FilterByFilter(r.Gcontext, r.FilterByList, r.FilterCustomizeFunc),
		QueryBySelect(r.Gcontext),
		QueryByPreload(r.PreloadList),
	).Table(r.ResourceTableName).First(r.ModelSerializer).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			r.HandleResponse(400, nil, err, gorm.ErrRecordNotFound)
			return

		}
		r.HandleResponse(400, nil, err, ErrLoadDataFailed)
		return
	}
	r.HandleResponse(200, r.ModelSerializer, nil, nil)
	return
}
