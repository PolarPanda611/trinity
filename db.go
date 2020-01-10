package trinity

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" //pg
	uuid "github.com/satori/go.uuid"
)

//InitDatabase create db connection
/**
 * initial db connection
 */
func (t *Trinity) InitDatabase() {
	var dbconnection string
	switch t.setting.Database.Type {
	case "mysql":
		var dbconn strings.Builder
		// 向builder中写入字符 / 字符串
		dbconn.Write([]byte(t.setting.Database.User))
		dbconn.WriteByte(':')
		dbconn.Write([]byte(t.setting.Database.Password))
		dbconn.Write([]byte("@/"))
		dbconn.Write([]byte(t.setting.Database.Name))
		dbconn.WriteByte('?')
		dbconn.Write([]byte(t.setting.Database.Option))
		dbconnection = dbconn.String()

		break
	case "postgres":
		dbconnection = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
			t.setting.Database.Host,
			t.setting.Database.Port,
			t.setting.Database.User,
			t.setting.Database.Password,
			t.setting.Database.Name,
			t.setting.Database.Option,
		)
		break
	}
	db, err := gorm.Open(t.setting.Database.Type, dbconnection)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return t.setting.Database.TablePrefix + defaultTableName
	}

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	db.LogMode(t.setting.Runtime.Debug)
	db.SetLogger(t.logger)
	db.SingularTable(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampAndUUIDForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	db.DB().SetMaxIdleConns(t.setting.Database.DbMaxIdleConn)
	db.DB().SetMaxOpenConns(t.setting.Database.DbMaxOpenConn)
	t.db = db

}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampAndUUIDForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		userID, ok := scope.Get("UserID")
		if !ok {
			userID = nil
		}
		nowTime := time.Now()
		if createTimeField, ok := scope.FieldByName("CreatedTime"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}
		if createUserIDField, ok := scope.FieldByName("CreateUserID"); ok {
			if createUserIDField.IsBlank {
				createUserIDField.Set(userID)
			}
		}
		if idField, ok := scope.FieldByName("ID"); ok {
			if idField.IsBlank {
				idField.Set(GenerateSnowFlakeID(int64(rand.Intn(100))))
			}
		}
		if modifyTimeField, ok := scope.FieldByName("UpdatedTime"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
		if updateUserIDField, ok := scope.FieldByName("UpdateUserID"); ok {
			if updateUserIDField.IsBlank {
				updateUserIDField.Set(userID)
			}
		}

		if updateDVersionField, ok := scope.FieldByName("DVersion"); ok {
			if updateDVersionField.IsBlank {
				updateDVersionField.Set(uuid.NewV4().String())
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		userID, ok := scope.Get("UserID")
		if !ok {
			userID = nil
		}
		var updateAttrs = map[string]interface{}{}
		if attrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
			updateAttrs = attrs.(map[string]interface{})
			updateAttrs["updated_time"] = time.Now()
			updateAttrs["update_user_id"] = userID
			updateAttrs["d_version"] = uuid.NewV4().String()
			scope.InstanceSet("gorm:update_attrs", updateAttrs)
		}
	}

}

// deleteCallback will set `DeletedOn` where deleting
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		userID, ok := scope.Get("UserID")
		if !ok {
			userID = nil
		}
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		deletedAtField, hasDeletedAtField := scope.FieldByName("deleted_time")
		deleteUserIDField, hasDeleteUserIDField := scope.FieldByName("DeleteUserID")
		dVersionField, hasDVersionField := scope.FieldByName("d_version")

		if !scope.Search.Unscoped && hasDeletedAtField && hasDVersionField && hasDeleteUserIDField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v,%v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedAtField.DBName),
				scope.AddToVars(time.Now()),
				scope.Quote(deleteUserIDField.DBName),
				scope.AddToVars(userID),
				scope.Quote(dVersionField.DBName),
				scope.AddToVars(uuid.NewV4().String()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
