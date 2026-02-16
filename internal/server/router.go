package server

import (
	"github.com/bdzhalalov/kolikosoft-trade/internal/item"
	"net/http"
)

func Router(itemHandler *item.Handler) http.Handler {
	rootRouter := http.NewServeMux()

	apiRouter := http.NewServeMux()

	item.RegisterRoutes(apiRouter, itemHandler)

	rootRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	return rootRouter
}
