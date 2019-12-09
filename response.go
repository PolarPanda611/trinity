package trinity

// ResponseData http response
type ResponseData struct {
	Status  int         // the http response status  to return
	Result  interface{} // the response data  if req success
	TraceID string
}

// Response handle trinity return value
func (v *ViewSetRunTime) Response() {
	var res ResponseData
	res.Status = v.Status
	res.TraceID = v.Gcontext.GetString("TraceID")
	if v.RealError != nil {
		// v.Cfg.Logger.LogWriter(v)
		v.Gcontext.Error(v.RealError)
		v.Gcontext.Error(v.UserError)
		res.Result = v.UserError.Error()
		v.Gcontext.AbortWithStatusJSON(v.Status, res)
	} else {
		res.Result = v.ResBody
		v.Gcontext.JSON(v.Status, res)
	}
	return

}

//HandleResponse handle response
func (v *ViewSetRunTime) HandleResponse(status int, payload interface{}, rerr error, uerr error) {
	v.mu.Lock()
	funcName, file, line := RecordErrorLevelTwo()
	v.FuncName = funcName
	v.File = file
	v.Line = line
	v.Status = status
	v.ResBody = payload
	v.RealError = rerr
	v.UserError = uerr
	v.mu.Unlock()
	v.Response()
}
