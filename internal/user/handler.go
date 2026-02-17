package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/render"
	"net/http"
	"strconv"
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

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	//TODO: add validation
	idStr := r.PathValue("id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || userId <= 0 {
		render.JSON(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var body WithdrawRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Amount <= 0 {
		render.JSON(w, "invalid amount: Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	requestId := r.Header.Get("Idempotency-Key")
	if requestId == "" {
		requestId, err = h.newRequestID()
		if err != nil {
			http.Error(w, "failed to generate request id", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Idempotency-Key", requestId)
	}

	dto := WithdrawBalanceRequestDTO{
		UserId:    userId,
		Amount:    body.Amount,
		RequestId: requestId,
	}

	res, e := h.service.WithdrawFromBalance(ctx, dto)
	if e != nil {
		render.JSON(w, e.Message, e.Code)
		return
	}

	render.JSON(w, res, http.StatusOK)
}

func (h *Handler) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) newRequestID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
