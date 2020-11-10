package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-exit:
			cancel()
		}
	}()

	config, err := parseConfig()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	container := newContainer(config)
	if err := container.Run(ctx).Wait(); err != nil {
		log.Fatal(err)
	}
}
