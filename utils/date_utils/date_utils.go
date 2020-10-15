package date_utils

import (
	"time"
)

func GetNowUnix() int64 {
	return time.Now().UTC().Unix()
}
