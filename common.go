package trinity

import (
	"time"
)

//Model common type
type Model struct {
	ID            uint       `json:"id"  gorm:"primary_key;"`
	Key           string     `json:"key" sql:"type:varchar(50);index" gorm:"unique;not null;default:'';"`
	CreatedTime   time.Time  `json:"created_time" sql:"index;"`
	CreateUser    *User      `json:"create_user" gorm:"AssociationForeignKey:CreateUserKey;ForeignKey:Key;"`
	CreateUserKey *string    `json:"create_user_key" gorm:"type:varchar(50);not null;default:'';" `
	UpdatedTime   time.Time  `json:"updated_time" sql:"index;"`
	UpdateUser    *User      `json:"update_user" gorm:"AssociationForeignKey:UpdateUserKey;ForeignKey:Key;"`
	UpdateUserKey *string    `json:"update_user_key" gorm:"type:varchar(50);not null;default:'';" `
	DeletedTime   *time.Time `json:"deleted_time" sql:"index;"`
	DeleteUser    *User      `json:"delete_user" gorm:"AssociationForeignKey:DeleteUserKey;ForeignKey:Key;"`
	DeleteUserKey *string    `json:"delete_user_key" gorm:"type:varchar(50);not null;default:'';" `
}

//Simpmodel common type
type Simpmodel struct {
	ID          uint      `json:"id"  gorm:"primary_key;"`
	CreatedTime time.Time `json:"created_time" sql:"index;"`
}

//Logmodel common type
type Logmodel struct {
	ID            uint      `json:"id,omitempty"  gorm:"primary_key;"`
	CreatedTime   time.Time `json:"created_time,omitempty" sql:"index;"`
	CreateUserKey *string   `json:"create_user_key,omitempty" gorm:"type:varchar(50);not null;default:'';" `
}

//Viewmodel for view type
type Viewmodel struct {
	ID            uint      `json:"id"  gorm:"primary_key;"`
	Key           string    `json:"key" sql:"type:varchar(50);index" gorm:"unique;not null;"`
	CreatedTime   time.Time `json:"created_time" sql:"index;"`
	CreateUserKey *string   `json:"create_user_key" gorm:"type:varchar(50);not null;default:'';" `
	// CreateUserUsername string     `json:"create_user_username"`
	// CreateUserRealname string     `json:"create_user_realname"`
	UpdatedTime   time.Time `json:"updated_time" sql:"index;"`
	UpdateUserKey *string   `json:"update_user_key" gorm:"type:varchar(50);not null;default:'';" `
	// UpdateUserUsername string     `json:"update_user_username"`
	// UpdateUserRealname string     `json:"update_user_realname"`
	DeletedTime   *time.Time `json:"deleted_time" sql:"index;"`
	DeleteUserKey *string    `json:"delete_user_key" gorm:"type:varchar(50);not null;default:'';" `
	// DeleteUserUsername string     `json:"delete_user_username"`
	// DeleteUserRealname string     `json:"delete_user_realname"`
}
