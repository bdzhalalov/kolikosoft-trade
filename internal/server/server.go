package server

import (
	"context"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/config"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/database"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/logger"
)

func Start(ctx context.Context, config *config.Config) {
	logger := logger.Logger(config)

	db := database.ConnectToDB(ctx, config)
	defer db.Close()

}
