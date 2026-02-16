package item

import "net/http"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
