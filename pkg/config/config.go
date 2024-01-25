package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostgresDSN       string `mapstructure:"DSN"`
	Env               string `mapstructure:"ENV"`
	HttpPort          string `mapstructure:"HTTP_PORT"`
	HttpHost          string `mapstructure:"HTTP_HOST"`
	AgeApiUrl         string `mapstructure:"AGE_API_URL"`
	GenderApiUrl      string `mapstructure:"GENDER_API_URL"`
	NationalityApiUrl string `mapstructure:"NATIONALITY_API_URL"`
}

/*// Get Postgres Url
func (c *Config) GetPsqlUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.PostgresUser, c.PostgresPass, c.PostgresHost, c.PostgresPort, c.PostgresDB, c.PostgresSSLMode)
}*/

// LoadConfig use viper lib to work with env
func LoadConfig(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv() // Allow automatically override values

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)

	return config, err
}
