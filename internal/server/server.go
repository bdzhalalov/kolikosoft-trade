package server

import (
	"context"
	"errors"
	"github.com/bdzhalalov/kolikosoft-trade/internal/item"
	"github.com/bdzhalalov/kolikosoft-trade/internal/user"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/cache"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/config"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/database"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/logger"
	"net/http"
	"time"
)

func Start(ctx context.Context, config *config.Config) {
	log := logger.Logger(config)

	db := database.ConnectToDB(ctx, config)

	httpClient := &http.Client{}
	c := cache.New()

	skinPortClient := item.NewSkinPortClient(httpClient, config.SkinPortBaseURL)
	itemService := item.NewService(skinPortClient, log, c)
	itemHandler := item.NewHandler(itemService)

	userRepo := user.NewRepository(db)
	userService := user.NewService(log, userRepo)
	userHandler := user.NewHandler(userService)

	router := Router(itemHandler, userHandler)

	apiServer := &http.Server{
		Addr:    config.Addr,
		Handler: router,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Infof("Running API server on port: %s", config.Addr)

		errCh <- apiServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Info("Shutting down API server...")
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("API server stopped unexpectedly")
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = apiServer.Shutdown(shutdownCtx)
	_ = db.Close()

	log.Info("API server shutdown complete")
}
