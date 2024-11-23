package common_test

import (
	"testing"
	"time"

	"github.com/delly/amartha/common"
	"github.com/stretchr/testify/assert"
)

func TestStartOfDay(t *testing.T) {
	t.Parallel()

	t.Run("should return start of day", func(t *testing.T) {
		t.Parallel()

		tm := time.Date(2021, 9, 1, 12, 30, 0, 0, time.UTC)
		expected := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)

		result := common.StartOfDay(tm)

		assert.Equal(t, expected, result)
	})
}

func TestEndOfDay(t *testing.T) {
	t.Parallel()

	t.Run("should return end of day", func(t *testing.T) {
		t.Parallel()

		tm := time.Date(2021, 9, 1, 12, 30, 0, 0, time.UTC)
		expected := time.Date(2021, 9, 1, 23, 59, 59, 0, time.UTC)

		result := common.EndOfDay(tm)

		assert.Equal(t, expected, result)
	})
}
