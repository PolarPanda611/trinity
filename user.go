package trinity

//User model User
type User struct {
	Model
	Username            string       `json:"username" gorm:"type:varchar(50);index;not null;"`             // login username /profile
	NameLocal           string       `json:"name_local"  gorm:"type:varchar(50);" `                        // local name
	NameEN              string       `json:"name_en"  gorm:"type:varchar(50);" `                           // EN name
	Email               string       `json:"email"  gorm:"type:varchar(50);" `                             // login email
	Phone               string       `json:"phone" gorm:"type:varchar(50);" `                              // login phone
	UserGroup           []Group      `json:"user_group" gorm:"many2many:user_group;"`                      // foreign key -->group
	UserPermission      []Permission `json:"user_permission" gorm:"many2many:user_permission;"`            // foreign key --->permission
	PreferenceLanguages string       `json:"preference_language" gorm:"type:varchar(50);default:'en-US'" ` // user preference language

}
