package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bdzhalalov/kolikosoft-trade/internal/user/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

var InsufficientFundsError = errors.New("insufficient funds")

func (r *Repository) GetUserById(ctx context.Context, userId int64) (domain.User, error) {
	const query = `SELECT id, balance FROM users WHERE id = $1`

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, userId).Scan(&u.Id, &u.Balance)
	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}

func (r *Repository) WithdrawFromUserBalance(
	ctx context.Context,
	userId int64,
	amount int64,
	requestId string,
) (domain.Withdrawal, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return domain.Withdrawal{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var w domain.Withdrawal
	row := tx.QueryRowContext(ctx, `
		SELECT user_id, amount, balance_before, balance_after, created_at
		FROM balance_withdrawals
		WHERE user_id = $1 AND request_id = $2
	`, userId, requestId)

	switch err := row.Scan(&w.UserId, &w.Amount, &w.BalanceBefore, &w.BalanceAfter, &w.CreatedAt); {
	case err == nil:
		if err := tx.Commit(); err != nil {
			return domain.Withdrawal{}, err
		}
		return w, nil
	case errors.Is(err, sql.ErrNoRows):
	default:
		return domain.Withdrawal{}, err
	}

	var balanceAfter int64
	err = tx.QueryRowContext(ctx, `
		UPDATE users
		SET balance = balance - $2
		WHERE id = $1 AND balance >= $2
		RETURNING balance
	`, userId, amount).Scan(&balanceAfter)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Withdrawal{}, InsufficientFundsError
		}
		return domain.Withdrawal{}, err
	}

	balanceBefore := balanceAfter + amount

	var res domain.Withdrawal
	err = tx.QueryRowContext(ctx, `
		INSERT INTO balance_withdrawals(user_id, request_id, amount, balance_before, balance_after)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id, amount, balance_before, balance_after, created_at
	`, userId, requestId, amount, balanceBefore, balanceAfter).Scan(
		&res.UserId,
		&res.Amount,
		&res.BalanceBefore,
		&res.BalanceAfter,
		&res.CreatedAt,
	)
	if err != nil {
		// если вдруг пришёл дубль request_id (гоночный retry) — можно вернуть существующую запись
		// но если у тебя request_id стабилен, сюда почти не попадёшь
		return domain.Withdrawal{}, err
	}

	if err := tx.Commit(); err != nil {
		return domain.Withdrawal{}, err
	}

	return res, nil
}

func (r *Repository) GetUserBalanceHistory(ctx context.Context, userId int64) ([]domain.Withdrawal, error) {
	const query = `
		SELECT user_id, amount, balance_before, balance_after, created_at
		FROM balance_withdrawals
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.Withdrawal

	for rows.Next() {
		var w domain.Withdrawal

		err := rows.Scan(
			&w.UserId,
			&w.Amount,
			&w.BalanceBefore,
			&w.BalanceAfter,
			&w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		history = append(history, w)
	}

	// ОБЯЗАТЕЛЬНО проверить ошибку rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
