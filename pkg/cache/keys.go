package cache

import (
	"time"
)

const (
	DefaultCacheTime time.Duration = 60 * 60 * time.Second
	Forever                        = 0
)

func GetConfigKey() string {
	return "usms_config:"
}
