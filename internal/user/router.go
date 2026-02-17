package user

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /users/{id}/balance/withdraw", h.Withdraw)
	mux.HandleFunc("GET /users/{id}/balance/history", h.GetBalanceHistory)
}
