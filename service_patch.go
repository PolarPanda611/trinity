package trinity

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"
)

// PatchHandler : List method
func PatchHandler(r *PatchMixin) {
	// if Callback BeforeGet registered , run before get callback
	if IsFuncInited(r.ViewSetRunTime.BeforePatch) {
		r.ViewSetRunTime.BeforePatch(r.ViewSetRunTime)
	}
	// if Callback Get registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.Patch) {
		r.ViewSetRunTime.Patch(r.ViewSetRunTime)
	}
	// if Callback AfterGet registered , run before get callback, if not , run default get callback
	if IsFuncInited(r.ViewSetRunTime.AfterPatch) {
		r.ViewSetRunTime.AfterPatch(r.ViewSetRunTime)
	}
	if r.ViewSetRunTime.RealError == nil {
		r.ViewSetRunTime.Response()
	}
}

// DefaultPatchCallback : PATCH method
func DefaultPatchCallback(r *ViewSetRunTime) {

	buf := make([]byte, 1024)
	n, _ := r.Gcontext.Request.Body.Read(buf)
	requestbodyMap := make(map[string]interface{})
	if err := json.Unmarshal(buf[0:n], &requestbodyMap); err != nil {
		r.HandleResponse(400, nil, err, ErrResolveDataFailed)
		return
	}
	id, err := strconv.ParseInt(r.Gcontext.Params.ByName("id"), 10, 64)
	if err != nil {
		r.HandleResponse(400, nil, err, ErrUpdateDataFailed)
		return
	}
	delete(requestbodyMap, "key")
	if r.EnableChangeLog {
		searchOldValue := reflect.New(reflect.ValueOf(r.ModelSerializer).Elem().Type()).Interface()

		if err := r.Db.Scopes(
			r.DBFilterBackend,
			FilterByParam(id),
			FilterByFilter(r.Gcontext, r.FilterByList, r.FilterCustomizeFunc),
			FilterBySearch(r.Gcontext, r.SearchingByList),
		).Table(r.ResourceTableName).First(searchOldValue).Error; err != nil {
			r.HandleResponse(400, nil, err, ErrUpdateDataFailed)
			return
		}
		v := reflect.ValueOf(searchOldValue).Elem()
		k := v.Type()
		oldDataMap := make(map[string]string)
		for i := 0; i < v.NumField(); i++ {
			key := gorm.ToColumnName(k.Field(i).Name)
			val := v.Field(i)
			oldDataMap[key] = fmt.Sprint(val)
			model, ok := val.Interface().(Model)
			if ok {
				oldDataMap["d_version"] = model.DVersion
			}
		}
		userID := r.Gcontext.GetInt64("UserID")
		for k, v := range requestbodyMap {
			oldValue := oldDataMap[k]
			newValue := fmt.Sprint(reflect.ValueOf(v))
			if oldValue == newValue {
				continue
			}
			changeLog := AppChangelog{
				Logmodel:    Logmodel{CreateUserID: userID},
				Resource:    r.ResourceTableName,
				Type:        "Update",
				Column:      k,
				OldValue:    oldValue,
				NewValue:    newValue,
				DVersion:    oldDataMap["d_version"],
				TraceID:     r.Gcontext.GetString("TraceID"),
				ResourceKey: r.Gcontext.Params.ByName("id"),
			}
			if err := r.Db.Create(&changeLog).Error; err != nil {
				r.HandleResponse(400, nil, err, ErrUpdateDataFailed)
				return
			}
		}
	}
	if r.EnableVersionControl {

	}
	updateQuery := r.Db.Set("gorm:save_associations", false).Scopes(
		r.DBFilterBackend,
		FilterByParam(id),
		FilterByFilter(r.Gcontext, r.FilterByList, r.FilterCustomizeFunc),
		FilterBySearch(r.Gcontext, r.SearchingByList),
		DataVersionFilter(requestbodyMap["d_version"], r.EnableDataVersion),
	).Table(r.ResourceTableName).Updates(requestbodyMap)

	if err := updateQuery.Error; err != nil {
		r.HandleResponse(400, nil, err, ErrUpdateDataFailed)
		return
	}
	if effectNum := updateQuery.RowsAffected; effectNum != 1 {
		r.HandleResponse(400, nil, ErrUpdateZeroAffected, ErrUpdateDataFailed)
		return
	}
	r.HandleResponse(200, "Updated Successfully", nil, nil)
	return

}
