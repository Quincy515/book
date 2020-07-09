package model

import (
	"blog/global"
	"blog/pkg/setting"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:id`
	CreatedBy  string `json:created_by`
	ModifiedBy string `json:modified_by`
	CreatedOn  uint32 `json:created_on`
	ModifiedOn uint32 `json:modified_on`
	DeletedOn  uint32 `json:deleted_on`
	IsDel      uint8  `json:is_del`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettings) (*gorm.DB, error) {
	db, err := gorm.Open(databaseSetting.DBType,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
			databaseSetting.Username,
			databaseSetting.Password,
			databaseSetting.Host,
			databaseSetting.DBName,
			databaseSetting.Charset,
			databaseSetting.ParseTime))
	if err != nil {
		return nil, err
	}
	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)
	// 注册回调行为
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
	return db, nil
}

// 新增行为的回调
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	// 通过调用 scope.FieldByName 方法获取是否包含所需的字段
	// 通过判断 Field.IsBlank 的值得知该字段的值是否为空
	// 若为空，则调用 Field.Set 方法给该字段设置值
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreateOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

// 更新行为的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	// 通过调用 scope.Get 来获取当前设置表示为 gorm:update_column 的字段属性
	// 若不存在，即没有自定义设置 update_column，则在更新回调内设置默认字段 ModifiedOn 的值为当前的时间戳
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 删除行为的回调
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		// 通过调用 scope.Get("gorm:delete_option") 来获取当前设置的标识 gorm:delete_option 的字段属性
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		// 判断是否存在 DeletedOn 和 IsDel 字段
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			now := time.Now().Unix()
			// 若存在则执行 UPDATE 进行软删除，修改 DeletedOn 和 IsDel 的值
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v, %v=%v%v%v",
				scope.QuotedTableName(), // 获取当前引用的表名
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExit(scope.CombinedConditionSql()), // scope.CombinedConditionSql 完成 SQL 语句的组装
				addExtraSpaceIfExit(extraOption),
			)).Exec()
		} else {
			// 否则执行 DELETE 操作进行硬删除
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExit(scope.CombinedConditionSql()),
				addExtraSpaceIfExit(extraOPtion),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExit(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
