package server

import (
	"github.com/bdzhalalov/kolikosoft-trade/internal/item"
	"github.com/bdzhalalov/kolikosoft-trade/internal/user"
	"net/http"
)

func Router(itemHandler *item.Handler, userHandler *user.Handler) http.Handler {
	rootRouter := http.NewServeMux()

	apiRouter := http.NewServeMux()

	item.RegisterRoutes(apiRouter, itemHandler)
	user.RegisterRoutes(apiRouter, userHandler)

	rootRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	return rootRouter
}
