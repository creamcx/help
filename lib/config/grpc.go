package config

import (
	"fmt"
	"net"
	"os"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type GrpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (*GrpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", grpcHostEnvName)

	}
	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", grpcPortEnvName)
	}
	return &GrpcConfig{
		host: host,
		port: port,
	}, nil
}

func (g *GrpcConfig) Address() string {
	return net.JoinHostPort(g.host, g.port)
}
