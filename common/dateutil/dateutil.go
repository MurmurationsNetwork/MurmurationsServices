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
