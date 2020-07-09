package global

import (
	"blog/pkg/logger"
	"blog/pkg/setting"
)

// 将配置信息和应用程序关联
var (
	ServerSetting   *setting.ServerSettings
	AppSetting      *setting.AppSettings
	DatabaseSetting *setting.DatabaseSettings
	Logger          *logger.Logger // 包全局变量中新增 Logger 对象，用于日志组件的初始化
)
