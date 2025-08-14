package smtp

import "time"

type Config struct {
	Host        string
	Port        int
	User        string
	Password    string
	Sender      string
	OtpDuration time.Duration
}
