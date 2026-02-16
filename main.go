package main

import (
	"context"
	"github.com/bdzhalalov/kolikosoft-trade/internal/server"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/config"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.InitConfig()

	server.Start(ctx, &cfg)
}
