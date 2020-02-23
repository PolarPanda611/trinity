package trinity

import "github.com/jinzhu/copier"

// PostHandler : List method
func PostHandler(r *PostMixin) {
	// if Callback BeforeGet registered , run before get callback
	if IsFuncInited(r.ViewSetRunTime.BeforePost) {
		r.ViewSetRunTime.BeforePost(r.ViewSetRunTime)
	}
	// if Callback Get registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.Post) {
		r.ViewSetRunTime.Post(r.ViewSetRunTime)
	}
	// if Callback AfterGet registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.AfterPost) {
		r.ViewSetRunTime.AfterPost(r.ViewSetRunTime)
	}
	if r.ViewSetRunTime.RealError == nil {
		r.ViewSetRunTime.Response()
	}
}

// DefaultPostCallback : Create method
func DefaultPostCallback(r *ViewSetRunTime) {
	if r.PostValidation != nil {
		if err := r.Gcontext.BindJSON(r.PostValidation); err != nil {
			r.HandleResponse(400, nil, err, ErrResolveDataFailed)
			return
		}
		copier.Copy(r.ModelSerializer, r.PostValidation)
	} else {
		if err := r.Gcontext.BindJSON(r.ModelSerializer); err != nil {
			r.HandleResponse(400, nil, err, ErrResolveDataFailed)
			return
		}
	}
	if err := r.Db.Set("gorm:save_associations", false).Table(r.ResourceTableName).Create(r.ModelSerializer).Error; err != nil {
		r.HandleResponse(400, nil, err, ErrCreateDataFailed)
		return
	}
	r.HandleResponse(201, r.ModelSerializer, nil, nil)
	return
}
