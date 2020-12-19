package dateutil

import (
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
