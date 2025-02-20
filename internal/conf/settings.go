package conf

import (
	"os"
)

type Config struct {
	FlaskMicroserviceURL string
}

func LoadConfig() *Config {
	return &Config{
		FlaskMicroserviceURL: getEnv("FLASK_MICROSERVICE_URL", "http://localhost:5000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
