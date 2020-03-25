package trinity

import (
	"github.com/jinzhu/gorm"
)

// Context  context impl
type Context struct {
	db      *gorm.DB
	logger  Logger
	setting ISetting
}

// GetDB get db instance
func (c *Context) GetDB(method string, traceID string, userName string) *gorm.DB {

	logger := &defaultLogger{
		ProjectName:    c.setting.GetProjectName(),
		ProjectVersion: c.setting.GetProjectVersion(),
		WebAppAddress:  c.setting.GetWebAppAddress(),
		WebAppPort:     c.setting.GetWebAppPort(),
	}
	logger.FormatLogger(method, traceID, userName)
	c.db.SetLogger(logger)
	return c.db
}
