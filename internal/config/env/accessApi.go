package env

import (
	"errors"
	"net"
	"os"
)

const (
	accessApiHost = "ACCESS_API_HOST"
	accessApiPort = "ACCESS_API_PORT"
)

type accessApiConfig struct {
	host string
	port string
}

func NewAccessApiClient() (*accessApiConfig, error) {
	host := os.Getenv(accessApiHost)

	port := os.Getenv(accessApiPort)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	return &accessApiConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *accessApiConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
