package consumer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

//go:generate mockgen -source=consumer.go -destination ./consumer_mocks.go -package consumer

type Reader interface {
	ReadMessage(ctx context.Context) ([]byte, error)
}

type Writer interface {
	WriteMessage(ctx context.Context, message []byte) error
}

type Decoder interface {
	Decode(b []byte) time.Time
}

type Consumer struct {
	writer  Writer
	reader  Reader
	decoder Decoder
}

func New(reader Reader, writer Writer, decoder Decoder) *Consumer {
	return &Consumer{
		reader:  reader,
		writer:  writer,
		decoder: decoder,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		date := c.decoder.Decode(message).Format(time.RFC3339)

		if err := c.writer.WriteMessage(ctx, []byte(date)); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
	}
}
