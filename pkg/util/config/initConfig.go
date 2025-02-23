package config

import "github.com/spf13/viper"

type Config struct {
	DbSource     string `mapstructure:"DB_SOURCE"`
	MigrationUrl string `mapstructure:"MIGRATION_URL"`
	SecretKey    string `mapstructure:"SECRET_KEY"`
}

func InitConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return config, err

}
