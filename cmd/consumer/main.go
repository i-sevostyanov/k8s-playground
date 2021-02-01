package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/kelseyhightower/envconfig"

	"github.com/i-sevostyanov/k8s-playground/internal/business/consumer"
	"github.com/i-sevostyanov/k8s-playground/internal/infrastructure/binary"
	"github.com/i-sevostyanov/k8s-playground/internal/infrastructure/kafka"
	"github.com/i-sevostyanov/k8s-playground/internal/infrastructure/metrics"
)

var errStopped = errors.New("stopped")

type Config struct {
	KafkaAddr   string `envconfig:"KAFKA_ADDR" default:"localhost:9092"`
	InputTopic  string `envconfig:"INPUT_TOPIC" default:"input"`
	OutputTopic string `envconfig:"OUTPUT_TOPIC" default:"output"`
	Listen      string `envconfig:"LISTEN" default:"localhost:8080"`
}

func main() {
	var cfg Config

	if err := envconfig.Process("consumer", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	reader := kafka.NewReader([]string{cfg.KafkaAddr}, cfg.InputTopic)
	writer := kafka.NewWriter(cfg.KafkaAddr, cfg.OutputTopic)
	encoder := binary.NewEncoder()
	aConsumer := consumer.New(reader, writer, encoder)

	gr, ctx := errgroup.WithContext(context.Background())

	gr.Go(func() error {
		return aConsumer.Run(ctx)
	})

	gr.Go(func() error {
		return metrics.Listen(ctx, cfg.Listen)
	})

	gr.Go(func() error {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(stop)

		select {
		case <-stop:
			log.Printf("Caught stop signal. Exiting...")
			return errStopped
		case <-ctx.Done():
			return nil
		}
	})

	log.Printf("Consumer started with config: %+v\n", cfg)

	if err := gr.Wait(); err != nil && err != errStopped {
		log.Printf("Error during wait: %v", err)
	}

	if err := reader.Close(); err != nil {
		log.Printf("Failed to close kafka reader: %v", err)
	}

	if err := writer.Close(); err != nil {
		log.Printf("Failed to close kafka writer: %v", err)
	}
}
