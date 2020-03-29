package trinity

import (
	"context"

	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/metadata"
)

// GRPCMethod grpc method
type GRPCMethod string

// ReqUserName request user name
type ReqUserName string

// TraceID current request trace id
type TraceID string

// UserRequestsCtx handle user ctx from context
type UserRequestsCtx interface {
	GetGRPCMethod() GRPCMethod
	GetTraceID() TraceID
	GetReqUserName() ReqUserName
}

// UserRequestsCtxImpl UserRequestsCtxImpl handler request ctx , extrace data from ctx
type UserRequestsCtxImpl struct {
	ctx         context.Context
	method      GRPCMethod
	reqUserName ReqUserName
	traceID     TraceID
}

// GetGRPCMethod get grpc method
func (u *UserRequestsCtxImpl) GetGRPCMethod() GRPCMethod { return u.method }

// GetTraceID get trace id
func (u *UserRequestsCtxImpl) GetTraceID() TraceID { return u.traceID }

// GetReqUserName get request user
func (u *UserRequestsCtxImpl) GetReqUserName() ReqUserName { return u.reqUserName }

// NewUserRequestsCtx new user ctx
func NewUserRequestsCtx(ctx context.Context) UserRequestsCtx {
	if ctx != nil {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			method := md[GRPCMethodKey][0]
			traceID := md[TraceIDKey][0]
			userName := md[ReqUserNameKey][0]
			return &UserRequestsCtxImpl{
				method:      GRPCMethod(method),
				reqUserName: ReqUserName(userName),
				traceID:     TraceID(traceID),
			}
		}
	}
	// no data carries in ctx
	return &UserRequestsCtxImpl{
		method:      "",
		reqUserName: "",
		traceID:     "",
	}
}

// Context interface to get Req Context
type Context interface {
	GetDB() *gorm.DB
	GetTXDB() *gorm.DB
	GetLogger() Logger
	GetRequest() interface{}
	GetResponse() interface{}
}

// ContextImpl  context impl
type ContextImpl struct {
	userRequestsCtx UserRequestsCtx
	db              *gorm.DB
	logger          Logger
	setting         ISetting
	request         interface{}
	response        interface{}
}

// NewContext new  ctx
func NewContext(db *gorm.DB, setting ISetting, userRequestsCtx UserRequestsCtx, logger Logger) Context {
	newDB := db.New()
	newDB.SetLogger(logger)
	newContext := &ContextImpl{
		userRequestsCtx: userRequestsCtx,
		logger:          logger,
		db:              newDB,
		setting:         setting,
	}
	return newContext
}

// GetDB get db instance
func (c *ContextImpl) GetDB() *gorm.DB {
	return c.db
}

// GetTXDB get db instance with transaction
func (c *ContextImpl) GetTXDB() *gorm.DB {
	c.db = c.db.Begin()
	return c.db
}

// GetLogger get current user name
func (c *ContextImpl) GetLogger() Logger {
	return c.logger
}

// GetRequest get user request data
func (c *ContextImpl) GetRequest() interface{} { return c.request }

// GetResponse get user response data
func (c *ContextImpl) GetResponse() interface{} { return c.response }
