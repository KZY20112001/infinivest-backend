package commons

import (
	"fmt"
	"time"
)

var (
	ValidFrequencies = map[string]struct{}{
		"daily":      {},
		"monthly":    {},
		"biweekly":   {},
		"weekly":     {},
		"quarterly":  {},
		"biannually": {},
		"annually":   {},
	}
)

func GetNextRebalanceTime(freq string) (time.Time, error) {
	nextRebalanceTime := time.Now()
	switch freq {
	case "daily":
		return nextRebalanceTime.AddDate(0, 0, 1), nil
	case "weekly":
		return nextRebalanceTime.AddDate(0, 0, 7), nil
	case "biweekly":
		return nextRebalanceTime.AddDate(0, 0, 14), nil
	case "monthly":
		return nextRebalanceTime.AddDate(0, 1, 0), nil
	case "quarterly":
		return nextRebalanceTime.AddDate(0, 3, 0), nil
	case "biannually":
		return nextRebalanceTime.AddDate(0, 6, 0), nil
	case "annually":
		return nextRebalanceTime.AddDate(1, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid rebalance frequency: %s", freq)
	}
}
