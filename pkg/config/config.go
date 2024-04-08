package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PostgresDSN       string
	Env               string
	HttpPort          string
	HttpHost          string
	AgeApiUrl         string
	GenderApiUrl      string
	NationalityApiUrl string
}

func (c *Config) GetHTTPPort() string {
	return c.HttpPort
}

func (c *Config) GetEnv() string {
	return c.Env
}

func (c *Config) GetAgeApiURL() string {
	return c.AgeApiUrl
}
func (c *Config) GetGenderApiURL() string {
	return c.GenderApiUrl
}
func (c *Config) GetNationalityApiURL() string {
	return c.NationalityApiUrl
}

func LoadEnv(filenames ...string) error {
	const op = "pkg.config.LoadEnv"
	err := godotenv.Load(filenames...)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	return nil
}

func LoadConfig() *Config {
	cfg := &Config{
		PostgresDSN:       "postgres://admin:qwerty@localhost:5432/human?sslmode=disable",
		Env:               "local",
		HttpHost:          "localhost",
		AgeApiUrl:         "https://api.agify.io/",
		GenderApiUrl:      "https://api.genderize.io/",
		NationalityApiUrl: "https://api.nationalize.io/",
	}

	postgresDsn := os.Getenv("DSN")
	env := os.Getenv("ENV")
	httpPort := os.Getenv("HTTP_PORT")
	httpHost := os.Getenv("HTTP_HOST")
	ageUrl := os.Getenv("AGE_API_URL")
	genderUrl := os.Getenv("GENDER_API_URL")
	nationalityUrl := os.Getenv("NATIONALITY_API_URL")

	if postgresDsn != "" {
		cfg.PostgresDSN = postgresDsn
	}
	if env != "" {
		cfg.Env = env
	}
	if httpPort != "" {
		cfg.HttpPort = httpPort
	}
	if httpHost != "" {
		cfg.HttpHost = httpHost
	}
	if ageUrl != "" {
		cfg.AgeApiUrl = ageUrl
	}
	if genderUrl != "" {
		cfg.GenderApiUrl = genderUrl
	}
	if nationalityUrl != "" {
		cfg.NationalityApiUrl = nationalityUrl
	}

	return cfg
}
