package trinity

import (
	"context"

	"github.com/jinzhu/gorm"
)

// GRPCMethod grpc method
type GRPCMethod string

// ReqUserName request user name
type ReqUserName string

// TraceID current request trace id
type TraceID string

// Context  context impl
type Context struct {
	db      *gorm.DB
	logger  Logger
	setting ISetting
	root    bool

	//traceid
	traceID TraceID

	//ReqUserName
	reqUserName ReqUserName
}

// NewContext new  ctx
func NewContext(ctx context.Context, db *gorm.DB, setting ISetting) *Context {
	method, traceID, reqUserName := GetLogFromMetaData(ctx)
	logger := &defaultLogger{
		ProjectName:    setting.GetProjectName(),
		ProjectVersion: setting.GetProjectVersion(),
		WebAppAddress:  setting.GetWebAppAddress(),
		WebAppPort:     setting.GetWebAppPort(),
	}
	logger.FormatLogger(method, traceID, reqUserName)
	newDB := db.New()
	newDB.SetLogger(logger)
	newContext := &Context{
		logger:      logger,
		db:          newDB,
		setting:     setting,
		traceID:     traceID,
		reqUserName: reqUserName,
	}
	return newContext
}

// GetDB get db instance
func (c *Context) GetDB() *gorm.DB {
	return c.db
}

// GetTraceID get current trace id
func (c *Context) GetTraceID() TraceID {

	return c.traceID
}

// GetReqUserName get current user name
func (c *Context) GetReqUserName() ReqUserName {
	return c.reqUserName
}
