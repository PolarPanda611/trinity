package trinity

//AppError model Error
type AppError struct {
	Logmodel
	TraceID  string `json:"trace_id"  gorm:"type:varchar(50);index;not null;"` //http seq number
	File     string `json:"file"  `
	Line     string `json:"line"  gorm:"type:varchar(50);"`
	FuncName string `json:"func_name"`
	Error    string `json:"error" `
}

// RecordError to record error
func (e *AppError) RecordError() {
	GlobalTrinity.db.Create(e)
}
