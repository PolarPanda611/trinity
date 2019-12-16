package trinity

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

//Permission model Role
type Permission struct {
	Model
	Code string `json:"code" gorm:"type:varchar(100);index;unique;not null;"`
	Name string `json:"name" gorm:"type:varchar(100);"`
}

// BeforeCreate hooks
func (p *Permission) BeforeCreate(scope *gorm.Scope) error {
	//add customize primary key
	scope.SetColumn("Key", uuid.NewV4().String())
	return nil
}

// CreateOrInitPermission create or init permission
func (p *Permission) CreateOrInitPermission(t *Trinity) {
	t.db.Where(Permission{Code: p.Code}).FirstOrCreate(p)
}

// PermissionViewSet hanlde router
func PermissionViewSet(c *gin.Context) {
	v := NewViewSet()
	v.HasAuthCtl = true
	v.NewRunTime(
		c,
		&Permission{},
		&Permission{},
		&[]Permission{},
	).ViewSetServe()

}
