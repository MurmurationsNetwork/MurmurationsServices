package backoff

import (
	"fmt"
	"time"

	backoff "github.com/cenkalti/backoff/v4"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
)

// NewBackoff is a utility function that implements an exponential backoff
// mechanism for retrying a given operation.
func NewBackoff(operation func() error, message string) error {
	b := backoff.NewExponentialBackOff()
	// 30 s -> 1 mins -> 2 mins -> 4 mins
	b.InitialInterval = 30 * time.Second
	b.RandomizationFactor = 0
	b.Multiplier = 2
	b.MaxInterval = 4 * time.Minute
	b.MaxElapsedTime = 15 * time.Minute

	err := backoff.RetryNotify(
		operation,
		b,
		func(err error, time time.Duration) {
			logger.Info(
				fmt.Sprintf(
					"%s, %s, retry in %0.f seconds \n",
					message,
					err,
					time.Seconds(),
				),
			)
		},
	)
	if err != nil {
		return err
	}

	return nil
}
