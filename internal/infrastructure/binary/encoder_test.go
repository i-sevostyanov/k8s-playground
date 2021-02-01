package binary_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/i-sevostyanov/k8s-playground/internal/infrastructure/binary"
)

func TestEncoder(t *testing.T) {
	t.Run("encode-decode works correct", func(t *testing.T) {
		expected := time.Now().UTC()
		encoder := binary.NewEncoder()
		bytes := encoder.Encode(expected)
		actual := encoder.Decode(bytes)

		assert.Equal(t, expected.Unix(), actual.Unix())
	})
}
