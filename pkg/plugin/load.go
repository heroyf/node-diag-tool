package plugin

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func DynamicLoadPlugins(path string, fileName string) {
	logrus.Debugf("read config: [%s].toml begin()", fileName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(path)
	viper.SetConfigName(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Debugf("read config failed, use default config instead. err: %v", err)
	} else {
		plugins := viper.Sub("plugins")
		EnabledPlugins = append(EnabledPlugins, plugins.GetStringSlice("enabledPlugins")...)
		DiabledPlugins = append(DiabledPlugins, plugins.GetStringSlice("disabledPlugins")...)
	}
}
