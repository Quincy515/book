package setting

import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper
}

// 初始化本项目配置的基础属性
func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")   // 设置配置文件名称 config
	vp.AddConfigPath("configs/") // 设置配置路径为相对路径 configs/
	vp.SetConfigType("yaml")     //  设置配置文件类型 yaml
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Setting{vp}, nil
}
