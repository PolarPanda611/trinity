package trinity

import (
	"github.com/PolarPanda611/gorm"
	"github.com/gin-gonic/gin"
)

//DefaultFilterBackend for Mixin
func DefaultFilterBackend(c *gin.Context, db *gorm.DB) *gorm.DB {
	// c.GetString("UserID")
	// c.GetStringSlice("UserPermission")
	return db
}
