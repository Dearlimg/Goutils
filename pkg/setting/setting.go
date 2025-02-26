package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

/*
使用 viper 进行配置文件的读取和热加载
配置热更新：开源库 github.com/fsnotify/fsnotify
*/

type Setting struct {
	vp  *viper.Viper
	all interface{} //用于存储配置文件中的所有配置信息
}

// NewSetting 初始化项目的基础属性
func NewSetting(configName, configType string, configPaths ...string) (*Setting, error) {
	//创建一个新的 viper 对象
	vp := viper.New()
	//设置配置文件的名称
	vp.SetConfigName(configName)
	//设置配置文件的类型
	vp.SetConfigType(configType)
	//设置配置文件的路径（可设置多个路径）
	for _, path := range configPaths {
		if path != "" {
			vp.AddConfigPath(path) //可设置多个配置路径，解决路径查找问题
		}
	}
	//加载配置文件
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}
	s := &Setting{vp: vp}
	//实时监控配置文件的变化
	s.vp.WatchConfig()
	//当配置变化之后调用的一个回调函数
	s.vp.OnConfigChange(func(in fsnotify.Event) {
		log.Println("更新配置")
		err := s.vp.Unmarshal(s.all)
		if err != nil {
			log.Fatalln("更新配置失败：" + err.Error())
		}
	})
	return s, nil
}

// BindAll 绑定配置文件
func (s *Setting) BindAll(v interface{}) error {
	//绑定
	err := s.vp.Unmarshal(v)
	if err != nil {
		return err
	}
	s.all = v
	return nil
}
