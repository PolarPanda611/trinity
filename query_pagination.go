package trinity

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//QueryByPagination handling pagination
func QueryByPagination(c *gin.Context, PageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		PageNumField := c.DefaultQuery("PageNum", "1")
		PageSizeField := c.DefaultQuery("PageSize", string(PageSize))
		var err error
		var PageNumFieldInt int
		var PageSizeFieldInt int
		if PageNumFieldInt, err = strconv.Atoi(PageNumField); err != nil || PageNumFieldInt < 0 {
			PageNumFieldInt = 0
		}
		PageNumFieldInt = PageNumFieldInt - 1
		if PageSizeFieldInt, err = strconv.Atoi(PageSizeField); err != nil {
			PageSizeFieldInt = PageSize
		}
		offset := PageNumFieldInt * PageSizeFieldInt
		limit := PageSizeFieldInt
		return db.Offset(offset).Limit(limit)

	}
}
