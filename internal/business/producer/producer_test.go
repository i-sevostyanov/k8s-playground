package producer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProducer_Run(t *testing.T) {
	t.Run("successful publish messages N times", func(t *testing.T) {
		const times = 5

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), times*time.Millisecond)
		defer cancel()

		message := []byte{1, 2, 3}
		encoder := NewMockEncoder(ctrl)
		writer := NewMockWriter(ctrl)

		encoder.EXPECT().Encode(gomock.Any()).Return(message).MaxTimes(times)
		writer.EXPECT().WriteMessage(ctx, message).Return(nil).MaxTimes(times)

		producer := New(writer, encoder)
		err := producer.Run(ctx, time.Millisecond)
		assert.NoError(t, err)
	})

	t.Run("returns error on write message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
		defer cancel()

		writeErr := errors.New("something went wrong")
		message := []byte{1, 2, 3}
		encoder := NewMockEncoder(ctrl)
		writer := NewMockWriter(ctrl)

		encoder.EXPECT().Encode(gomock.Any()).Return(message)
		writer.EXPECT().WriteMessage(ctx, message).Return(writeErr)

		producer := New(writer, encoder)
		err := producer.Run(ctx, time.Millisecond)
		assert.Equal(t, writeErr, errors.Unwrap(err))
	})

	t.Run("returns nil when the context is done", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		encoder := NewMockEncoder(ctrl)
		writer := NewMockWriter(ctrl)

		producer := New(writer, encoder)
		err := producer.Run(ctx, time.Second)
		assert.NoError(t, err)
	})
}
