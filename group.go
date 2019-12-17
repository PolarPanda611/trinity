package trinity

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

//Group model Group
type Group struct {
	Model
	Name        string       `json:"name" gorm:"type:varchar(50);unique;not null;"`
	Description string       `json:"description" gorm:"type:varchar(100);"`
	Permissions []Permission `json:"permissions" gorm:"many2many:group_permission;AssociationForeignkey:Key;ForeignKey:Key;"`
	PGroup      *Group       `json:"p_group"  gorm:"AssociationForeignKey:PKey;Foreignkey:Key;"`
	PKey        string       `json:"p_key"`
	Roles       []Role       `json:"roles" gorm:"many2many:group_role;AssociationForeignkey:Key;ForeignKey:Key;"`
}

// BeforeCreate hooks
func (group *Group) BeforeCreate(scope *gorm.Scope) error {
	//add customize primary key
	scope.SetColumn("Key", uuid.NewV4().String())
	return nil
}
