package cache

import (
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsFresh(t *testing.T) {
	t.Run("should be fresh for new instance", func(t *testing.T) {
		c := New(123, time.Now)
		assert.Equal(t, false, c.IsFresh())
	})

	t.Run("should be fresh if time if not expired", func(t *testing.T) {
		calls := 0
		var maxAge int64 = 10
		now := time.Now()

		mockNow := func() time.Time {
			if calls == 0 {
				calls += 1
				return now
			}

			return now.Add(time.Duration(maxAge-1) * time.Second)
		}

		c := New(maxAge, mockNow)
		c.Update([]model.Redirect{})
		assert.Exactly(t, true, c.IsFresh())
	})

	t.Run("should not be fresh after time is expired", func(t *testing.T) {
		calls := 0
		var maxAge int64 = 10
		now := time.Now()
		mockNow := func() time.Time {
			if calls == 0 {
				calls += 1
				return now.Add(time.Duration(-(maxAge + 123)) * time.Second) // for Update set time in the pas
			}

			return now
		}

		c := New(maxAge, mockNow)
		c.Update([]model.Redirect{})
		assert.Exactly(t, false, c.IsFresh())
	})

	t.Run("should handle cache lock", func(t *testing.T) {
		c := New(10, time.Now)
		// first lock - should success
		lr := c.Lock()
		assert.Equal(t, true, lr)

		// second lock - should return false
		lr = c.Lock()
		assert.Equal(t, false, lr)

		c.Unlock()

		// lock after unlock - should success
		lr = c.Lock()
		assert.Equal(t, true, lr)
	})
}
