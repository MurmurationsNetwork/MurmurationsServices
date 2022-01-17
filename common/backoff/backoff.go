package backoff

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
)

func NewBackoff(operation func() error, message string) error {
	b := backoff.NewExponentialBackOff()
	// 30 s -> 1 mins -> 2 mins -> 4 mins
	b.InitialInterval = 30 * time.Second
	b.RandomizationFactor = 0
	b.Multiplier = 2
	b.MaxInterval = 4 * time.Minute
	b.MaxElapsedTime = 15 * time.Minute

	err := backoff.RetryNotify(operation, b, func(err error, time time.Duration) {
		logger.Info(fmt.Sprintf("%s, %s, retry in %0.f seconds \n", message, err, time.Seconds()))
	})
	if err != nil {
		return err
	}

	return nil
}
