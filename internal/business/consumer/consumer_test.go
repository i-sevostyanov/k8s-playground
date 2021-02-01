package consumer

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestConsumer(t *testing.T) {
	t.Run("successful read-write messages N times", func(t *testing.T) {
		const times = 5

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		now := time.Now().UTC()
		date := now.Format(time.RFC3339)
		message := []byte{1, 2, 3}

		reader := NewMockReader(ctrl)
		writer := NewMockWriter(ctrl)
		decoder := NewMockDecoder(ctrl)

		gomock.InOrder(
			reader.EXPECT().ReadMessage(ctx).Return(message, nil).Times(times),
			reader.EXPECT().ReadMessage(ctx).Return(nil, io.EOF),
		)

		decoder.EXPECT().Decode(message).Return(now).Times(times)
		writer.EXPECT().WriteMessage(ctx, []byte(date)).Return(nil).Times(times)

		consumer := New(reader, writer, decoder)
		err := consumer.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("returns error on read message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		reader := NewMockReader(ctrl)
		writer := NewMockWriter(ctrl)
		decoder := NewMockDecoder(ctrl)

		readErr := errors.New("something went wrong")
		reader.EXPECT().ReadMessage(ctx).Return(nil, readErr)

		consumer := New(reader, writer, decoder)
		err := consumer.Run(ctx)
		assert.Equal(t, readErr, errors.Unwrap(err))
	})

	t.Run("returns nil on EOF", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		reader := NewMockReader(ctrl)
		writer := NewMockWriter(ctrl)
		decoder := NewMockDecoder(ctrl)

		reader.EXPECT().ReadMessage(ctx).Return(nil, io.EOF)

		consumer := New(reader, writer, decoder)
		err := consumer.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("returns error on write message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		writeErr := errors.New("something went wrong")
		now := time.Now().UTC()
		date := now.Format(time.RFC3339)
		message := []byte{1, 2, 3}

		reader := NewMockReader(ctrl)
		writer := NewMockWriter(ctrl)
		decoder := NewMockDecoder(ctrl)

		reader.EXPECT().ReadMessage(ctx).Return(message, nil)
		decoder.EXPECT().Decode(message).Return(now)
		writer.EXPECT().WriteMessage(ctx, []byte(date)).Return(writeErr)

		consumer := New(reader, writer, decoder)
		err := consumer.Run(ctx)
		assert.Equal(t, writeErr, errors.Unwrap(err))
	})
}
