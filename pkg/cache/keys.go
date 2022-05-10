package cache

import (
	"time"
)

const (
	DefaultMinutesTime time.Duration = 60 * 60 * time.Second
	DefaultHourTime                  = 1 * time.Hour
	Forever                          = 0
)

func GetConfigKey() string {
	return "usms_config:"
}
