package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/config"
	_ "github.com/lib/pq"
	"time"
)

func ConnectToDB(ctx context.Context, config *config.Config) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUser, config.DbPass, config.DbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		panic(err)
	}

	return db
}
