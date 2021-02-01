package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	conn *kafka.Reader
}

func NewReader(brokers []string, topic string) *Reader {
	conn := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	return &Reader{
		conn: conn,
	}
}

func (r *Reader) Close() error {
	return r.conn.Close()
}

func (r *Reader) ReadMessage(ctx context.Context) ([]byte, error) {
	message, err := r.conn.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	return message.Value, nil
}
