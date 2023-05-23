package dateutil

import (
	"fmt"
	"time"
)

func GetNowUnix() int64 {
	return time.Now().UTC().Unix()
}

func GetZeroValueUnix() int64 {
	return time.Time{}.UTC().Unix()
}

func NowSubtract(duration time.Duration) int64 {
	return time.Now().UTC().Add(duration * -1).Unix()
}

func FormatSeconds(seconds int64) string {
	duration := time.Duration(seconds) * time.Second

	years := duration / (365 * 24 * time.Hour)
	duration -= years * (365 * 24 * time.Hour)

	months := duration / (30 * 24 * time.Hour)
	duration -= months * (30 * 24 * time.Hour)

	weeks := duration / (7 * 24 * time.Hour)
	duration -= weeks * (7 * 24 * time.Hour)

	days := duration / (24 * time.Hour)
	duration -= days * (24 * time.Hour)

	hours := duration / time.Hour
	duration -= hours * time.Hour

	minutes := duration / time.Minute
	duration -= minutes * time.Minute

	seconds = int64(duration / time.Second)

	result := ""
	if years > 0 {
		result += fmt.Sprintf("%d years", years)
	}
	if months > 0 {
		if len(result) > 0 {
			result += " "
		}
		result += fmt.Sprintf("%d months", months)
	}
	if weeks > 0 {
		if len(result) > 0 {
			result += " "
		}
		result += fmt.Sprintf("%d weeks", weeks)
	}
	if days > 0 {
		if len(result) > 0 {
			result += " "
		}
		result += fmt.Sprintf("%d days", days)
	}
	if hours > 0 {
		if len(result) > 0 {
			result += " "
		}
		result += fmt.Sprintf("%d hours", hours)
	}
	if minutes > 0 {
		if len(result) > 0 {
			result += " "
		}
		result += fmt.Sprintf("%d minutes", minutes)
	}
	if seconds > 0 {
		if len(result) > 0 {
			result += " "
		}
		result += fmt.Sprintf("%d seconds", seconds)
	}

	return result
}
