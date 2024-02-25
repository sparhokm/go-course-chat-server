package env

import (
	"fmt"
	"os"
)

var envList = map[string]string{
	"host":     "POSTGRES_HOST",
	"port":     "POSTGRES_PORT",
	"name":     "POSTGRES_DB",
	"user":     "POSTGRES_USER",
	"password": "POSTGRES_PASSWORD",
}

type pgConfig struct {
	dsn string
}

func NewPGConfig() (*pgConfig, error) {
	env := make(map[string]string, len(envList))

	var envVal string
	for i, v := range envList {
		envVal = os.Getenv(v)
		if len(envVal) == 0 {
			return nil, fmt.Errorf("%s host not found", i)
		}
		env[i] = envVal
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		env["host"],
		env["port"],
		env["name"],
		env["user"],
		env["password"],
	)

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
