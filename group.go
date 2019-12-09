package trinity

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

//Group model Group
type Group struct {
	Model
	Name            string       `json:"name" gorm:"type:varchar(50);index;unique;not null;"`
	Description     string       `json:"description" gorm:"type:varchar(100);index;"`
	GroupPermission []Permission `json:"group_permission" gorm:"many2many:group_permission;AssociationForeignkey:Key;ForeignKey:Key;"`
}

// BeforeCreate hooks
func (group *Group) BeforeCreate(scope *gorm.Scope) error {
	//add customize primary key
	scope.SetColumn("Key", uuid.NewV4().String())
	return nil
}

// GroupViewSet hanlde router
func GroupViewSet(c *gin.Context) {
	v := NewViewSet()
	v.HasAuthCtl = true
	v.NewRunTime(
		c,
		&Group{},
		&Group{},
		&[]Group{},
	).ViewSetServe()

}
