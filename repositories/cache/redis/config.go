package redis

import "time"

type Config struct {
	Host        string
	Port        int
	Password    string
	Database    int
	OtpDuration time.Duration
}
