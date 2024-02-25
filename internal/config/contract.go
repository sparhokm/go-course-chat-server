package config

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}
