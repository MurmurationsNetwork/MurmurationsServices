package retry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
)

// dummy function that will fail a few times before succeeding.
func testFunc() func() error {
	var callCount = 0
	return func() error {
		callCount++
		if callCount < 3 {
			return fmt.Errorf("call number %d failed", callCount)
		}
		// success on the 3rd call.
		return nil
	}
}

func TestRetry(t *testing.T) {
	retryFunc := testFunc()

	err := retry.Do(
		retryFunc,
		retry.WithInitialBackoff(1*time.Millisecond),
		retry.WithMaxBackoff(1*time.Second),
		retry.WithMultiplier(2),
		retry.WithRandomizationFactor(0.5),
	)

	require.NoError(t, err)
}
