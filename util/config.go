package util

import "github.com/spf13/viper"

type Config struct {
	FromEmail    string `mapstructure:"FROM_EMAIL"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	ToEmail      string `mapstructure:"TO_EMAIL"`
	ApiAddress   string `mapstructure:"API_ADDRESS"`
	ApiAuth      string `mapstructure:"API_AUTH"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
	

}