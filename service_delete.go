package trinity

import "strconv"

// DeleteHandler : DELETE method
func DeleteHandler(r *DeleteMixin) {
	// if Callback BeforeDelete registered , run before Delete callback
	if IsFuncInited(r.ViewSetRunTime.BeforeDelete) {
		r.ViewSetRunTime.BeforeDelete(r.ViewSetRunTime)
	}
	// if Callback Delete registered , run before Delete callback, if not , run default Delete callback
	if IsFuncInited(r.ViewSetRunTime.Delete) {
		r.ViewSetRunTime.Delete(r.ViewSetRunTime)
	}

	// if Callback AfterDelete registered , run before Delete callback, if not , run default Delete callback
	if IsFuncInited(r.ViewSetRunTime.AfterDelete) {
		r.ViewSetRunTime.AfterDelete(r.ViewSetRunTime)
	}
	r.ViewSetRunTime.Response()
}

// DefaultDeleteCallback default delete callback
func DefaultDeleteCallback(r *ViewSetRunTime) {
	id, err := strconv.ParseInt(r.Gcontext.Params.ByName("id"), 10, 64)
	if err != nil {
		r.HandleResponse(400, nil, err, ErrDeleteDataFailed)
		return
	}
	if err := r.Db.Scopes(
		r.DBFilterBackend,
		FilterByParam(id),
		FilterByFilter(r.Gcontext, r.FilterByList, r.FilterCustomizeFunc),
		FilterBySearch(r.Gcontext, r.SearchingByList),
	).Table(r.ResourceTableName).Delete(r.ModelSerializer).Error; err != nil {
		r.HandleResponse(400, nil, err, ErrDeleteDataFailed)
		return
	}
	r.HandleResponse(200, "Delete Successfully", nil, nil)
	return
}
