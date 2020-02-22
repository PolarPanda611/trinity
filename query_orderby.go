package trinity

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//QueryByOrdering handle ordering
func QueryByOrdering(c *gin.Context, OrderingByList map[string]bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//Searching Section : keywords : OrderingBy , separate by comma
		//example : OrderingBy=-id,user,-name , - means desc , default asc
		OrderByField := c.Query("OrderingBy")
		if OrderByField == "" {
			return db.Order("id asc")
		}
		ordercondition := ""
		for _, orderField := range strings.Split(OrderByField, ",") {
			if len(strings.Split(orderField, "-")) > 1 {
				if _, ok := OrderingByList[strings.Split(orderField, "-")[1]]; ok {
					ordercondition += strings.Split(orderField, "-")[1] + " desc ,"
				}
			} else {
				if _, ok := OrderingByList[orderField]; ok {
					ordercondition += orderField + " asc ,"
				}
			}
		}
		return db.Order(strings.TrimSuffix(ordercondition, ","))
	}
}
