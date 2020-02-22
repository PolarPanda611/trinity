package trinity

import "strconv"

// DeleteResource : DELETE method
func DeleteResource(r *DeleteMixin) {
	id, err := strconv.ParseInt(r.ViewSetRunTime.Gcontext.Params.ByName("id"), 10, 64)
	if err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, ErrDeleteDataFailed)
		return
	}
	if err := r.ViewSetRunTime.Db.Scopes(
		r.ViewSetRunTime.DBFilterBackend,
		FilterByParam(id),
		FilterByFilter(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.FilterByList, r.ViewSetRunTime.FilterCustomizeFunc),
		FilterBySearch(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.SearchingByList),
	).Table(r.ViewSetRunTime.ResourceTableName).Delete(r.ViewSetRunTime.ModelSerializer).Error; err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, ErrDeleteDataFailed)
		return
	}
	r.ViewSetRunTime.HandleResponse(200, "Delete Successfully", nil, nil)
	return
}
