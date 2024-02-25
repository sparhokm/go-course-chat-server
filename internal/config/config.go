package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/sparhokm/go-course-ms-chat-server/internal/config/env"
)

type Config struct {
	GRPCConfig GRPCConfig
	PGConfig   PGConfig
}

func MustLoad() *Config {
	path := fetchConfigPath()
	err := godotenv.Load(path)

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	return &Config{GRPCConfig: grpcConfig, PGConfig: pgConfig}
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", ".env", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
