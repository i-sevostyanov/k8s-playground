package producer

import (
	"context"
	"fmt"
	"time"
)

//go:generate mockgen -source=producer.go -destination ./producer_mocks.go -package producer

type Writer interface {
	WriteMessage(ctx context.Context, message []byte) error
}

type Encoder interface {
	Encode(t time.Time) []byte
}

type Producer struct {
	writer  Writer
	encoder Encoder
}

func New(writer Writer, encoder Encoder) *Producer {
	return &Producer{
		writer:  writer,
		encoder: encoder,
	}
}

func (p *Producer) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-ticker.C:
			message := p.encoder.Encode(t)

			if err := p.writer.WriteMessage(ctx, message); err != nil {
				return fmt.Errorf("failed to write message: %w", err)
			}
		}
	}
}
