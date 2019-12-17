package trinity

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

//AppError model Error
type AppError struct {
	Logmodel
	TraceID  string `json:"trace_id"  gorm:"type:varchar(50);index;not null;"` //http seq number
	File     string `json:"file"  `
	Line     string `json:"line"  gorm:"type:varchar(50);"`
	FuncName string `json:"func_name"`
	Error    string `json:"error" `
}

// BeforeCreate hooks
func (e *AppError) BeforeCreate(scope *gorm.Scope) error {
	//add customize primary key
	scope.SetColumn("Key", uuid.NewV4().String())
	return nil
}

// RecordError to record error
func (e *AppError) RecordError() {
	GlobalTrinity.db.Create(e)
}

//AppErrorViewSet for app error http handle
func AppErrorViewSet(c *gin.Context) {
	v := NewViewSet()
	v.HasAuthCtl = true
	v.FilterByList = []string{"trace_id"}
	v.NewRunTime(
		c,
		&AppError{},
		&AppError{},
		&[]AppError{},
	).ViewSetServe()
}
