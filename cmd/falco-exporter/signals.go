package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func withSignals(ctx context.Context) context.Context {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			return
		case s := <-sigCh:
			switch s {
			case os.Interrupt:
				log.Println("received SIGINT, shutting down")
			case syscall.SIGTERM:
				log.Println("received SIGTERM, shutting down")
			}
			return
		}
	}()
	return ctx
}
