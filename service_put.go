package trinity

import ()

// PutHandler : List method
func PutHandler(r *PutMixin) {
	// if Callback BeforeGet registered , run before get callback
	if IsFuncInited(r.ViewSetRunTime.BeforePut) {
		r.ViewSetRunTime.BeforePut(r.ViewSetRunTime)
	}
	// if Callback Get registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.Put) {
		r.ViewSetRunTime.Put(r.ViewSetRunTime)
	}
	// if Callback AfterGet registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.AfterPut) {
		r.ViewSetRunTime.AfterPut(r.ViewSetRunTime)
	}
	r.ViewSetRunTime.Response()
}

// DefaultPutCallback : PATCH method
func DefaultPutCallback(r *ViewSetRunTime) {
	r.HandleResponse(200, "PutMixin", nil, nil)
	return
}
