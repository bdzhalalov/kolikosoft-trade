package domain

import "time"

type User struct {
	Id      int64
	Balance int64
}

type Withdrawal struct {
	UserId        int64
	Amount        int64
	BalanceBefore int64
	BalanceAfter  int64
	CreatedAt     time.Time
}
