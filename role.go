package trinity

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

//Role model Role
type Role struct {
	Model
	Name        string       `json:"name" gorm:"type:varchar(50);index;unique;not null;"`
	Description string       `json:"description" gorm:"type:varchar(100);index;"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permission;AssociationForeignkey:Key;ForeignKey:Key;"`
}

// BeforeCreate hooks
func (r *Role) BeforeCreate(scope *gorm.Scope) error {
	//add customize primary key
	scope.SetColumn("Key", uuid.NewV4().String())
	return nil
}
