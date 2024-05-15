package initialize

import (
	"github.com/spf13/viper"
	"mxshop_srvs/usr_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {

	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFileName := "usr_srv/config.pro.yaml"
	if debug {
		configFileName = "usr_srv/config.dev.yaml"
	}
	v := viper.New()
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}
