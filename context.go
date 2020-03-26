package trinity

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
)

// Context  context impl
type Context struct {
	db      *gorm.DB
	logger  Logger
	setting ISetting
	root    bool

	//method
	method string

	//traceid
	traceID string

	//username
	userName string
}

// Init init ctx
func (c *Context) Init(ctx context.Context) *Context {
	method, traceID, userName := GetLogFromMetaData(ctx)
	logger := &defaultLogger{
		ProjectName:    c.setting.GetProjectName(),
		ProjectVersion: c.setting.GetProjectVersion(),
		WebAppAddress:  c.setting.GetWebAppAddress(),
		WebAppPort:     c.setting.GetWebAppPort(),
	}
	logger.FormatLogger(method, traceID, userName)
	newDB := c.db.New()
	newDB.SetLogger(logger)
	newContext := &Context{
		logger:   logger,
		db:       newDB,
		setting:  c.setting,
		root:     false,
		method:   method,
		traceID:  traceID,
		userName: userName,
	}
	return newContext
}

// GetDB get db instance
func (c *Context) GetDB() (*gorm.DB, error) {
	if c.root {
		return nil, errors.New("need init context first ")
	}
	return c.db, nil
}

// GetMethod get current Method
func (c *Context) GetMethod() (string, error) {
	if c.root {
		return "", errors.New("need init context first ")
	}
	return c.method, nil
}

// GetTraceID get current trace id
func (c *Context) GetTraceID() (string, error) {
	if c.root {
		return "", errors.New("need init context first ")
	}
	return c.traceID, nil
}

// GetUserName get current user name
func (c *Context) GetUserName() (string, error) {
	if c.root {
		return "", errors.New("need init context first ")
	}
	return c.userName, nil
}
