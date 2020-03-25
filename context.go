package trinity

import "github.com/jinzhu/gorm"

// Context  context impl
type Context struct {
	db     *gorm.DB
	logger Logger
}

// GetDB get db instance
func (c *Context) GetDB() *gorm.DB {
	return c.db
}

// GetLogger get logger instance
func (c *Context) GetLogger() Logger {
	return c.logger
}

// GetLogger get logger instance
func (c *Context) clone() *Context {
	cClone := &Context{
		db:     c.db,
		logger: c.logger.clone(),
	}
	return cClone
}
