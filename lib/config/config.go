package config

import (
	"github.com/spf13/viper"
	"lib/file"
)

// Config ...
type Config struct {
	path string
}

// Init ...
func Init(cfgPath string) error {
	c := Config{
		path: cfgPath,
	}
	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	currentPath, err := file.CurrentExecPath()
	if err != nil {
		return err
	}

	viper.AddConfigPath(currentPath + "/conf")

	if c.path != "" {
		viper.AddConfigPath(c.path)
	}

	//如果没有指定配置文件，则解析默认的配置文件

	var configFile string = "config_prod"

	dev, err := file.PathExists(currentPath + "/.dev")
	if err != nil {
		return err
	}
	if dev == true {
		configFile = "config_dev"
	}

	test, err := file.PathExists(currentPath + "/.test")
	if err != nil {
		return err
	}
	if test == true {
		configFile = "config_test"
	}

	prod, err := file.PathExists(currentPath + "/.prod")
	if err != nil {
		return err
	}
	if prod == true {
		configFile = "config_prod"
	}

	viper.SetConfigName(configFile)

	// 设置配置文件格式为YAML
	viper.SetConfigType("yaml")

	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// // 监听配置文件是否改变,用于热更新
// func (c *Config) watchConfig() {
//     viper.WatchConfig()
//     viper.OnConfigChange(func(e fsnotify.Event) {
//         fmt.Printf("Config file changed: %s\n", e.Name)
//     })
// }
