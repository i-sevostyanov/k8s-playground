package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
	conn *kafka.Writer
}

func NewWriter(addr, topic string) *Writer {
	conn := &kafka.Writer{
		Addr:         kafka.TCP(addr),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
	}

	return &Writer{
		conn: conn,
	}
}

func (w *Writer) Close() error {
	return w.conn.Close()
}

func (w *Writer) WriteMessage(ctx context.Context, message []byte) error {
	return w.conn.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
}
