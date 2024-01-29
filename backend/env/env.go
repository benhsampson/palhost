package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DBConnectionString string
	DBName             string
}

func (e *EnvConfig) Load(envFile string) {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal(err)
	}
	e.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	e.DBName = os.Getenv("DB_NAME")
}

func NewEnvConfig(envFile string) *EnvConfig {
	e := EnvConfig{}
	e.Load(envFile)
	return &e
}
