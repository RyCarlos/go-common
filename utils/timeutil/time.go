package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

func ParseDuration(durationStr string) (time.Duration, error) {
	// 支持 "7d", "24h", "30m" 等格式
	if len(durationStr) == 0 {
		return 0, fmt.Errorf("empty duration string")
	}

	unit := durationStr[len(durationStr)-1]
	valueStr := durationStr[:len(durationStr)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	switch unit {
	case 'd':
		return time.Duration(value) * 24 * time.Hour, nil
	case 'h':
		return time.Duration(value) * time.Hour, nil
	case 'm':
		return time.Duration(value) * time.Minute, nil
	case 's':
		return time.Duration(value) * time.Second, nil
	default:
		return 0, fmt.Errorf("unknown duration unit: %c", unit)
	}
}
