package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/bdzhalalov/kolikosoft-trade/internal/user/domain"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	logger   = logrus.New()
	mockUser = domain.User{
		Id:      1,
		Balance: 300,
	}
	history       = make([]domain.Withdrawal, 0, 2)
	mockRequestId = "request-id"
)

type RepositoryMock struct {
	err error
}

func (r *RepositoryMock) GetUserById(_ context.Context, userId int64) (domain.User, error) {
	if userId != 1 {
		return domain.User{}, sql.ErrNoRows
	}

	return mockUser, nil
}

func (r *RepositoryMock) WithdrawFromUserBalance(
	_ context.Context,
	userId int64,
	amount int64,
	requestId string,
) (domain.Withdrawal, error) {
	if r.err != nil {
		return domain.Withdrawal{}, r.err
	}

	if requestId == mockRequestId {
		return history[0], nil
	}

	if amount > mockUser.Balance {
		return domain.Withdrawal{}, InsufficientFundsError
	}
	withdrawal := domain.Withdrawal{
		UserId:        userId,
		Amount:        amount,
		BalanceBefore: mockUser.Balance,
		BalanceAfter:  mockUser.Balance - amount,
		CreatedAt:     time.Now(),
	}

	mockUser.Balance = mockUser.Balance - amount

	history = append(history, withdrawal)

	return withdrawal, nil
}

func (r *RepositoryMock) GetUserBalanceHistory(_ context.Context, userId int64) ([]domain.Withdrawal, error) {
	if r.err != nil {
		return []domain.Withdrawal{}, r.err
	}
	return history, nil
}

func seedBalanceHistory() {
	for _, _ = range history {
		history = append(history, domain.Withdrawal{
			UserId:        mockUser.Id,
			Amount:        50,
			BalanceBefore: mockUser.Balance,
			BalanceAfter:  mockUser.Balance - 50,
			CreatedAt:     time.Now(),
		})
	}
}

func TestWithdrawFromBalanceOk(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{},
		logger:     logger,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/1/balance/withdraw", strings.NewReader(`{"amount": 50}`))
	req.SetPathValue("id", "1")

	handler.Withdraw(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if len(history) != 1 {
		t.Fatalf("expected 1 history entity, got %d", len(history))
	}

	if history[0].Amount != 50 {
		t.Fatalf("expected amount to be 50, got %d", history[0].Amount)
	}

	if mockUser.Balance != history[0].BalanceAfter {
		t.Fatalf("expected balance to be %d, got %d", mockUser.Balance, history[0].BalanceAfter)
	}
}

func TestWithdrawFromBalanceIdempotency(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{},
		logger:     logger,
	}

	handler := NewHandler(svc)

	firstRec := httptest.NewRecorder()
	firstReq := httptest.NewRequest(http.MethodPost, "/api/v1/users/1/balance/withdraw", strings.NewReader(`{"amount": 50}`))
	firstReq.SetPathValue("id", "1")
	firstReq.Header.Set("Idempotency-Key", mockRequestId)

	handler.Withdraw(firstRec, firstReq)

	var bodyFirstResponse WithdrawBalanceResponseDTO

	err := json.NewDecoder(firstRec.Body).Decode(&bodyFirstResponse)
	if err != nil {
		t.Fatal(err)
	}

	if firstRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", firstRec.Code)
	}

	secondRec := httptest.NewRecorder()
	secondReq := httptest.NewRequest(http.MethodPost, "/api/v1/users/1/balance/withdraw", strings.NewReader(`{"amount": 40}`))
	secondReq.SetPathValue("id", "1")
	secondReq.Header.Set("Idempotency-Key", mockRequestId)

	handler.Withdraw(secondRec, secondReq)

	var bodySecondResponse WithdrawBalanceResponseDTO
	e := json.NewDecoder(secondRec.Body).Decode(&bodySecondResponse)
	if e != nil {
		t.Fatal(e)
	}

	if secondRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", secondRec.Code)
	}

	if len(history) != 1 {
		t.Fatalf("expected 1 history entity, got %d", len(history))
	}

	if bodyFirstResponse.BalanceAfter != bodySecondResponse.BalanceAfter {
		t.Fatalf(
			"Balance from idempotent request not equal to balance from first request. Expected %d, got %d",
			bodyFirstResponse.BalanceAfter, bodySecondResponse.BalanceAfter,
		)
	}
}

func TestWithdrawFromUnexistingUserBalance(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{},
		logger:     logger,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/2/balance/withdraw", strings.NewReader(`{"amount": 50}`))
	req.SetPathValue("id", "2")

	handler.Withdraw(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestWithdrawFromBalanceWithInvalidAmount(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{},
		logger:     logger,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/1/balance/withdraw", strings.NewReader(`{"amount": -40}`))
	req.SetPathValue("id", "1")

	handler.Withdraw(rec, req)

	var response string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	if response != "invalid amount: Amount must be greater than 0" {
		t.Fatalf("Expected %s, got %s", "invalid amount: Amount must be greater than 0", response)
	}

}

func TestWithdrawFromBalanceInsufficientFunds(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{},
		logger:     logger,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/1/balance/withdraw", strings.NewReader(`{"amount": 500}`))
	req.SetPathValue("id", "1")

	handler.Withdraw(rec, req)

	var response string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	if response != "insufficient funds" {
		t.Fatalf("expected %s, got %s", "insufficient funds", response)
	}
}

func TestWithdrawFromBalanceWithInvalidUserId(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{},
		logger:     logger,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/test/balance/withdraw", strings.NewReader(`{"amount": 40}`))
	req.SetPathValue("id", "test")

	handler.Withdraw(rec, req)

	var response string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	if response != "invalid user id" {
		t.Fatalf("expected %s, got %s", "invalid user id", response)
	}
}

func TestWithdrawFromBalanceErrorFromRepository(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{err: errors.New("some error from repository")},
		logger:     logger,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/1/balance/withdraw", strings.NewReader(`{"amount": 40}`))
	req.SetPathValue("id", "1")

	handler.Withdraw(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestGetBalanceHistoryOk(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{},
		logger:     logger,
	}

	seedBalanceHistory()

	handler := NewHandler(svc)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/balance/history", nil)
	req.SetPathValue("id", "1")

	handler.GetBalanceHistory(rec, req)

	var response []BalanceHistoryResponseDTO
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if len(response) != len(history) {
		t.Fatalf("expected %d history entity, got %d", len(history), len(response))
	}
}

func TestGetBalanceHistoryForUnexisitingUser(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{},
		logger:     logger,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/2/balance/history", nil)
	req.SetPathValue("id", "2")

	handler.GetBalanceHistory(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetBalanceHistoryWithInvalidUserId(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{},
		logger:     logger,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/test/balance/history", nil)
	req.SetPathValue("id", "test")

	handler.GetBalanceHistory(rec, req)

	var response string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	if response != "invalid user id" {
		t.Fatalf("expected %s, got %s", "invalid user id", response)
	}
}

func TestGetBalanceHistoryErrorFromRepository(t *testing.T) {
	svc := &Service{
		repository: &RepositoryMock{err: errors.New("some error from repository")},
		logger:     logger,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1/balance/history", nil)
	req.SetPathValue("id", "1")

	handler.GetBalanceHistory(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
