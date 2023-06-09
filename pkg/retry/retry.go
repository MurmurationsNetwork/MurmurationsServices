package retry

import (
	"fmt"
	"time"

	backoff "github.com/cenkalti/backoff/v4"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
)

// Default values for retry settings.
const (
	DefaultInitialBackoff      = 30 * time.Second
	DefaultMaxBackoff          = 4 * time.Minute
	DefaultMultiplier          = 2
	DefaultRandomizationFactor = 0.15
)

// Retry struct holds configuration for retry behavior.
type Retry struct {
	// Default retry backoff interval.
	InitialBackoff time.Duration
	// Maximum retry backoff interval.
	MaxBackoff time.Duration
	// Default backoff constant.
	Multiplier float64
	// Randomize the backoff interval by constant.
	RandomizationFactor float64
}

// Option is a function that modifies a Retry.
type Option func(*Retry)

// WithInitialBackoff sets the initial backoff value.
func WithInitialBackoff(d time.Duration) Option {
	return func(r *Retry) {
		r.InitialBackoff = d
	}
}

// WithMaxBackoff sets the maximum backoff value.
func WithMaxBackoff(d time.Duration) Option {
	return func(r *Retry) {
		r.MaxBackoff = d
	}
}

// WithMultiplier sets the multiplier for a Retry.
func WithMultiplier(m float64) Option {
	return func(r *Retry) {
		r.Multiplier = m
	}
}

// WithRandomizationFactor sets the randomization factor value.
func WithRandomizationFactor(randomizationFactor float64) Option {
	return func(r *Retry) {
		r.RandomizationFactor = randomizationFactor
	}
}

// Do executes the provided function with retry logic.
func Do(
	fn func() error,
	message string,
	opts ...Option,
) error {
	r := &Retry{
		InitialBackoff:      DefaultInitialBackoff,
		MaxBackoff:          DefaultMaxBackoff,
		Multiplier:          DefaultMultiplier,
		RandomizationFactor: DefaultRandomizationFactor,
	}

	// Apply provided options.
	for _, opt := range opts {
		opt(r)
	}

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = r.InitialBackoff
	b.MaxInterval = r.MaxBackoff
	b.Multiplier = r.Multiplier
	b.RandomizationFactor = r.RandomizationFactor

	return backoff.RetryNotify(
		fn,
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
}
