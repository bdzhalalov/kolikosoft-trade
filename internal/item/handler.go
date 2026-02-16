package item

import (
	"context"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/render"
	"net/http"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.service.GetItems(ctx)
	if err != nil {
		render.JSON(w, err.Message, err.Code)
	}

	render.JSON(w, items, http.StatusOK)
}
