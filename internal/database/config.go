package database

import "time"

type Config struct {
	Host            string
	Port            string
	User            string
	Pass            string
	Name            string
	SslMode         string
	MaxConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}
