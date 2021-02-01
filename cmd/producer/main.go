package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/kelseyhightower/envconfig"

	"github.com/i-sevostyanov/k8s-playground/internal/business/producer"
	"github.com/i-sevostyanov/k8s-playground/internal/infrastructure/binary"
	"github.com/i-sevostyanov/k8s-playground/internal/infrastructure/kafka"
	"github.com/i-sevostyanov/k8s-playground/internal/infrastructure/metrics"
)

var errStopped = errors.New("stopped")

type Config struct {
	KafkaAddr string        `envconfig:"KAFKA_ADDR" default:"localhost:9092"`
	Topic     string        `envconfig:"TOPIC" default:"input"`
	Interval  time.Duration `envconfig:"INTERVAL" default:"5s"`
	Listen    string        `envconfig:"LISTEN" default:"localhost:8080"`
}

func main() {
	var cfg Config

	if err := envconfig.Process("producer", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	writer := kafka.NewWriter(cfg.KafkaAddr, cfg.Topic)
	encoder := binary.NewEncoder()
	aProducer := producer.New(writer, encoder)

	gr, ctx := errgroup.WithContext(context.Background())

	gr.Go(func() error {
		return aProducer.Run(ctx, cfg.Interval)
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

	log.Printf("Producer started with config: %+v\n", cfg)

	if err := gr.Wait(); err != nil && err != errStopped {
		log.Printf("Error during wait: %v", err)
	}

	if err := writer.Close(); err != nil {
		log.Printf("Failed to close kafka writer: %v", err)
	}
}
