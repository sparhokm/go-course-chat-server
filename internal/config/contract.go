package config

type GRPCConfig interface {
	Address() string
}

type SwaggerConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type AccessApiConfig interface {
	Address() string
}
