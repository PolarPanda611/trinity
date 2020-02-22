package trinity

import "github.com/jinzhu/copier"

// CreateResource : Create method
func CreateResource(r *PostMixin) {
	if r.ViewSetRunTime.PostValidation != nil {
		if err := r.ViewSetRunTime.Gcontext.BindJSON(r.ViewSetRunTime.PostValidation); err != nil {
			r.ViewSetRunTime.HandleResponse(400, nil, err, ErrResolveDataFailed)
			return
		}
		copier.Copy(r.ViewSetRunTime.ModelSerializer, r.ViewSetRunTime.PostValidation)
	} else {
		if err := r.ViewSetRunTime.Gcontext.BindJSON(r.ViewSetRunTime.ModelSerializer); err != nil {
			r.ViewSetRunTime.HandleResponse(400, nil, err, ErrResolveDataFailed)
			return
		}
	}
	if err := r.ViewSetRunTime.Db.Set("gorm:save_associations", false).Table(r.ViewSetRunTime.ResourceTableName).Create(r.ViewSetRunTime.ModelSerializer).Error; err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, ErrCreateDataFailed)
		return
	}
	r.ViewSetRunTime.HandleResponse(201, r.ViewSetRunTime.ModelSerializer, nil, nil)
	return
}
