package otel

import "fmt"

type Config struct {
	ServiceName          string
	OtlpGrpcHost         string
	OtlpGrpcPort         int
	OtlpSecureConnection bool
}

func (cfg *Config) OtlpGrpcEndpoint() string {
	return fmt.Sprintf("%s:%d", cfg.OtlpGrpcHost, cfg.OtlpGrpcPort)
}
