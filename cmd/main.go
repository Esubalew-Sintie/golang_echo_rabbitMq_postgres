package main

import (
	"context"
	"os"
	"os/signal"
	"payment-gateway/initiator"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	initiator.Initialize(ctx)
}
