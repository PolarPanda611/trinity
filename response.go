package trinity

import (
	"encoding/json"
	"runtime"
	"strconv"
)

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
	res.TraceID = v.TraceID
	if v.RealError != nil {
		userKey := v.Gcontext.GetString("UserID")
		// v.Cfg.Logger.LogWriter(v)
		e := AppError{
			Logmodel: Logmodel{CreateUserKey: &userKey},
			TraceID:  v.Gcontext.GetString("TraceID"),
			File:     v.File,
			Line:     strconv.Itoa(v.Line),
			FuncName: runtime.FuncForPC(v.FuncName).Name(),
			Error:    v.RealError.Error()}
		b, _ := json.Marshal(e)
		v.Gcontext.Set("ErrorDetail", string(b))
		e.RecordError()

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
