package trinity

// RetrieveMixin for Get request
type RetrieveMixin struct {
	ViewSetRunTime *ViewSetRunTime
}

// GetMixin for Get request
type GetMixin struct {
	ViewSetRunTime *ViewSetRunTime
}

// PostMixin for Get request
type PostMixin struct {
	ViewSetRunTime *ViewSetRunTime
}

// PatchMixin for Get request
type PatchMixin struct {
	ViewSetRunTime *ViewSetRunTime
}

// PutMixin for Get request
type PutMixin struct {
	ViewSetRunTime *ViewSetRunTime
}

// DeleteMixin for Get request
type DeleteMixin struct {
	ViewSetRunTime *ViewSetRunTime
}

// UnknownMixin for Get request
type UnknownMixin struct {
	ViewSetRunTime *ViewSetRunTime
}

// Handler for handle retrieve request
func (r *RetrieveMixin) Handler() {
	GetResourceByid(r)
	return
}

// Handler for handle retrieve request
func (r *GetMixin) Handler() {
	GetResourceList(r)
	return
}

// Handler for handle post request
func (r *PostMixin) Handler() {
	CreateResource(r)
	return
}

// Handler for handle patch request
func (r *PatchMixin) Handler() {
	PatchResource(r)
	return
}

// Handler for handle put request
func (r *PutMixin) Handler() {
	r.ViewSetRunTime.HandleResponse(200, "PutMixin", nil, nil)
	return
}

// Handler for handle delete request
func (r *DeleteMixin) Handler() {
	DeleteResource(r)
	return
}

// Handler for handle Unknown request
func (r *UnknownMixin) Handler() {
	r.ViewSetRunTime.HandleResponse(405, nil, ErrUnknownService, ErrUnknownService)
	return
}

//ReqMixinHandler for handle mixin
type ReqMixinHandler interface {
	Handler()
}
