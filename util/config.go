package util

import (
	"errors"
	"github.com/spf13/viper"
	"path"
)

func InitConfig(inst interface{}, confPath string) error{
	vp := viper.New()
	vp.AddConfigPath(path.Dir(confPath))
	vp.SetConfigFile(path.Base(confPath))
	//Global.WatchConfig()
	//Global.OnConfigChange(func(e fsnotify.Event) {
	//    fmt.Println("配置发生变更：", e.Name)
	//})
	if err := vp.ReadInConfig(); err != nil {
		return errors.New("viper read err: " + err.Error())
	}
	if err := vp.Unmarshal(inst); err != nil {
		return errors.New("viper unmarshal err: " + err.Error())
	}
	return nil
}
