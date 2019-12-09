package trinity

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

//User model User
type User struct {
	Model
	Username            string       `json:"username" gorm:"type:varchar(50);index;unique;not null;"`                                    // login username /profile
	NameLocal           string       `json:"name_local"`                                                                                 // local name
	NameEN              string       `json:"name_en"`                                                                                    // EN name
	Email               string       `json:"email"`                                                                                      // login email
	Phone               string       `json:"phone" gorm:"type:varchar(50);" `                                                            // login phone
	UserGroup           []Group      `json:"user_group" gorm:"many2many:user_group;AssociationForeignkey:Key;ForeignKey:Key;"`           // foreign key -->group
	UserPermission      []Permission `json:"user_permission" gorm:"many2many:user_permission;AssociationForeignkey:Key;ForeignKey:Key;"` // foreign key --->permission
	PreferenceLanguages string       `json:"preference_language" gorm:"type:language;default:'en-US'" `                                  // user preference language

}

// BeforeCreate hooks
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	//add customize primary key
	scope.SetColumn("Key", uuid.NewV4().String())
	return nil
}

// UserViewSet hanlde router
func UserViewSet(c *gin.Context) {

	v := NewViewSet()
	v.HasAuthCtl = true
	v.NewRunTime(
		c,
		&User{},
		&User{},
		&[]User{},
	).ViewSetServe()
}
