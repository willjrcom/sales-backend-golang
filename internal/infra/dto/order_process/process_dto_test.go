package processdto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

func TestCalculateElapsedDuration(t *testing.T) {
	now := time.Date(2024, time.January, 1, 10, 0, 0, 0, time.UTC)

	t.Run("returns elapsed time for started process", func(t *testing.T) {
		startedAt := now.Add(-5 * time.Minute)
		process := &orderprocessentity.OrderProcess{
			OrderProcessCommonAttributes: orderprocessentity.OrderProcessCommonAttributes{
				Status: orderprocessentity.ProcessStatusStarted,
			},
			OrderProcessTimeLogs: orderprocessentity.OrderProcessTimeLogs{
				StartedAt: &startedAt,
			},
		}

		duration := calculateElapsedDuration(process, now)

		assert.Equal(t, 5*time.Minute, duration)
	})

	t.Run("includes accumulated duration when process continued", func(t *testing.T) {
		continuedAt := now.Add(-30 * time.Second)
		process := &orderprocessentity.OrderProcess{
			OrderProcessCommonAttributes: orderprocessentity.OrderProcessCommonAttributes{
				Status: orderprocessentity.ProcessStatusContinued,
			},
			OrderProcessTimeLogs: orderprocessentity.OrderProcessTimeLogs{
				Duration:    2 * time.Minute,
				ContinuedAt: &continuedAt,
			},
		}

		duration := calculateElapsedDuration(process, now)

		assert.Equal(t, 2*time.Minute+30*time.Second, duration)
	})

	t.Run("returns stored duration when process not running", func(t *testing.T) {
		process := &orderprocessentity.OrderProcess{
			OrderProcessCommonAttributes: orderprocessentity.OrderProcessCommonAttributes{
				Status: orderprocessentity.ProcessStatusPaused,
			},
			OrderProcessTimeLogs: orderprocessentity.OrderProcessTimeLogs{
				Duration: 4 * time.Minute,
			},
		}

		duration := calculateElapsedDuration(process, now)

		assert.Equal(t, 4*time.Minute, duration)
	})

	t.Run("ignores future reference times", func(t *testing.T) {
		startedAt := now.Add(5 * time.Minute)
		process := &orderprocessentity.OrderProcess{
			OrderProcessCommonAttributes: orderprocessentity.OrderProcessCommonAttributes{
				Status: orderprocessentity.ProcessStatusStarted,
			},
			OrderProcessTimeLogs: orderprocessentity.OrderProcessTimeLogs{
				StartedAt: &startedAt,
				Duration:  time.Minute,
			},
		}

		duration := calculateElapsedDuration(process, now)

		assert.Equal(t, time.Minute, duration)
	})
}

func TestFormatDurationToClock(t *testing.T) {
	testCases := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "minutes and seconds",
			duration: time.Minute + 5*time.Second,
			expected: "01:05",
		},
		{
			name:     "multi hour duration keeps minutes growing",
			duration: 125 * time.Minute,
			expected: "125:00",
		},
		{
			name:     "negative duration coerced to zero",
			duration: -10 * time.Second,
			expected: "00:00",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, formatDurationToClock(tc.duration))
		})
	}
}
