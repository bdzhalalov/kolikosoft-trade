package user

import "time"

type WithdrawBalanceRequestDTO struct {
	UserId    int64
	Amount    int64
	RequestId string
}

type WithdrawRequestBody struct {
	Amount int64 `json:"amount"`
}

type WithdrawBalanceResponseDTO struct {
	UserId        int64     `json:"user_id"`
	Amount        int64     `json:"amount"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}

type BalanceHistoryResponseDTO struct {
	UserId        int64     `json:"user_id"`
	Amount        int64     `json:"amount"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}
