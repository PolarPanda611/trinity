package trinity

import (
	"errors"
)

// RetrieveMixin for Get request
type RetrieveMixin struct {
	*ViewSetRunTime
}

// GetMixin for Get request
type GetMixin struct {
	*ViewSetRunTime
}

// PostMixin for Get request
type PostMixin struct {
	*ViewSetRunTime
}

// PatchMixin for Get request
type PatchMixin struct {
	*ViewSetRunTime
}

// PutMixin for Get request
type PutMixin struct {
	*ViewSetRunTime
}

// DeleteMixin for Get request
type DeleteMixin struct {
	*ViewSetRunTime
}

// UnknownMixin for Get request
type UnknownMixin struct {
	*ViewSetRunTime
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
	r.ViewSetRunTime.HandleResponse(405, nil, errors.New("The HyperText Transfer Protocol (HTTP) 405 Method Not Allowed response status code indicates that the request method is known by the server but is not supported by the target resource"), errors.New("The HyperText Transfer Protocol (HTTP) 405 Method Not Allowed response status code indicates that the request method is known by the server but is not supported by the target resource"))
	return
}

//ReqMixinHandler for handle mixin
type ReqMixinHandler interface {
	Handler()
}
