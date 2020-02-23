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
	RetriveHandler(r)
	return
}

// Handler for handle retrieve request
func (r *GetMixin) Handler() {
	GetHandler(r)
	return
}

// Handler for handle post request
func (r *PostMixin) Handler() {
	PostHandler(r)
	return
}

// Handler for handle patch request
func (r *PatchMixin) Handler() {
	PatchHandler(r)
	return
}

// Handler for handle put request
func (r *PutMixin) Handler() {
	PutHandler(r)
}

// Handler for handle delete request
func (r *DeleteMixin) Handler() {
	DeleteHandler(r)
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
