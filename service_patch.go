package trinity

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"
)

// PatchResource : PATCH method
func PatchResource(r *PatchMixin) {

	buf := make([]byte, 1024)
	n, _ := r.ViewSetRunTime.Gcontext.Request.Body.Read(buf)
	requestbodyMap := make(map[string]interface{})
	if err := json.Unmarshal(buf[0:n], &requestbodyMap); err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, ErrResolveDataFailed)
		return
	}
	id, err := strconv.ParseInt(r.ViewSetRunTime.Gcontext.Params.ByName("id"), 10, 64)
	if err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, ErrUpdateDataFailed)
		return
	}
	delete(requestbodyMap, "key")
	if r.ViewSetRunTime.EnableChangeLog {
		searchOldValue := reflect.New(reflect.ValueOf(r.ViewSetRunTime.ModelSerializer).Elem().Type()).Interface()

		if err := r.ViewSetRunTime.Db.Scopes(
			r.ViewSetRunTime.DBFilterBackend,
			FilterByParam(id),
			FilterByFilter(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.FilterByList, r.ViewSetRunTime.FilterCustomizeFunc),
			FilterBySearch(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.SearchingByList),
		).Table(r.ViewSetRunTime.ResourceTableName).First(searchOldValue).Error; err != nil {
			r.ViewSetRunTime.HandleResponse(400, nil, err, ErrUpdateDataFailed)
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
		userID := r.ViewSetRunTime.Gcontext.GetInt64("UserID")
		for k, v := range requestbodyMap {
			oldValue := oldDataMap[k]
			newValue := fmt.Sprint(reflect.ValueOf(v))
			if oldValue == newValue {
				continue
			}
			changeLog := AppChangelog{
				Logmodel:    Logmodel{CreateUserID: userID},
				Resource:    r.ViewSetRunTime.ResourceTableName,
				Type:        "Update",
				Column:      k,
				OldValue:    oldValue,
				NewValue:    newValue,
				DVersion:    oldDataMap["d_version"],
				TraceID:     r.ViewSetRunTime.Gcontext.GetString("TraceID"),
				ResourceKey: r.ViewSetRunTime.Gcontext.Params.ByName("id"),
			}
			if err := r.ViewSetRunTime.Db.Create(&changeLog).Error; err != nil {
				r.ViewSetRunTime.HandleResponse(400, nil, err, ErrUpdateDataFailed)
				return
			}
		}
	}
	if r.ViewSetRunTime.EnableVersionControl {

	}
	updateQuery := r.ViewSetRunTime.Db.Set("gorm:save_associations", false).Scopes(
		r.ViewSetRunTime.DBFilterBackend,
		FilterByParam(id),
		FilterByFilter(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.FilterByList, r.ViewSetRunTime.FilterCustomizeFunc),
		FilterBySearch(r.ViewSetRunTime.Gcontext, r.ViewSetRunTime.SearchingByList),
		DataVersionFilter(requestbodyMap["d_version"], r.ViewSetRunTime.EnableDataVersion),
	).Table(r.ViewSetRunTime.ResourceTableName).Updates(requestbodyMap)

	if err := updateQuery.Error; err != nil {
		r.ViewSetRunTime.HandleResponse(400, nil, err, ErrUpdateDataFailed)
		return
	}
	if effectNum := updateQuery.RowsAffected; effectNum != 1 {
		r.ViewSetRunTime.HandleResponse(400, nil, ErrUpdateZeroAffected, ErrUpdateDataFailed)
		return
	}
	r.ViewSetRunTime.HandleResponse(200, "Updated Successfully", nil, nil)
	return

}
