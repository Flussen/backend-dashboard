package config

type Config struct {
	User      string
	Password  string
	Host      string
	Port      string
	Database  string
	SecretKey string
}

var AppConfig *Config
