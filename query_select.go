package trinity

import (
	"encoding/json"
	"strings"

	"github.com/PolarPanda611/gorm"
	"github.com/gin-gonic/gin"
)

//QueryBySelect handling select
func QueryBySelect(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		SelectByField := c.Query("SelectBy")
		if SelectByField == "" {
			return db
		}
		strings.Split(SelectByField, ",")
		return db.Select(strings.Split(SelectByField, ","))

	}
}

//HandlestructBySelect handling select
func HandlestructBySelect(c *gin.Context, resource interface{}) (int, interface{}, error) {
	SelectByField := c.Query("SelectBy")
	if SelectByField == "" {
		return 200, resource, nil
	}
	m, _ := json.Marshal(&resource)
	// decode it back to get a map
	var a interface{}
	json.Unmarshal(m, &a)
	b := a.(map[string]interface{})
	for k := range b {
		if !stringInSlice(k, strings.Split(SelectByField, ",")) {
			delete(b, k)
		}
	}
	json.Marshal(b)
	return 200, b, nil

}
