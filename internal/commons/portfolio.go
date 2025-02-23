package commons

import (
	"fmt"
	"time"
)

var (
	RebalancingThresholds = map[string]float64{
		"daily":      10.0, // ±10%
		"weekly":     7.0,  // ±7%
		"biweekly":   5.0,  // ±5%
		"monthly":    3.0,  // ±3%
		"quarterly":  2.0,  // ±2%
		"biannually": 1.0,  // ±1%
		"annually":   0.5,  // ±0.5%
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
