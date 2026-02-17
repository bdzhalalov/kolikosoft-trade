package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bdzhalalov/kolikosoft-trade/internal/user/domain"
	customError "github.com/bdzhalalov/kolikosoft-trade/pkg/error"
	"github.com/sirupsen/logrus"
)

type RepositoryInterface interface {
	GetUserById(ctx context.Context, userId int64) (domain.User, error)
	WithdrawFromUserBalance(ctx context.Context, userId int64, amount int64, requestId string) (domain.Withdrawal, error)
}

type Service struct {
	logger     *logrus.Logger
	repository RepositoryInterface
}

func NewService(logger *logrus.Logger, repository RepositoryInterface) *Service {
	return &Service{
		logger,
		repository,
	}
}

func (s *Service) WithdrawFromBalance(
	ctx context.Context,
	data WithdrawBalanceRequestDTO,
) (WithdrawBalanceResponseDTO, *customError.BaseError) {
	_, err := s.repository.GetUserById(ctx, data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WithdrawBalanceResponseDTO{}, (&customError.NotFoundError{}).New("User not found")
		}

		s.logger.Errorf("Error while getting user by ID: %s", err)

		return WithdrawBalanceResponseDTO{}, (&customError.InternalServerError{}).New()
	}

	res, err := s.repository.WithdrawFromUserBalance(ctx, data.UserId, data.Amount, data.RequestId)
	if err != nil {
		if errors.Is(err, InsufficientFundsError) {
			return WithdrawBalanceResponseDTO{}, (&customError.BadRequestError{}).New("insufficient funds")
		}
		s.logger.Errorf("Error while withdrawal from user balance by ID: %s", err)
		return WithdrawBalanceResponseDTO{}, (&customError.InternalServerError{}).New()
	}

	dto := WithdrawBalanceResponseDTO{
		UserId:        res.UserId,
		Amount:        res.Amount,
		BalanceBefore: res.BalanceBefore,
		BalanceAfter:  res.BalanceAfter,
		CreatedAt:     res.CreatedAt,
	}

	return dto, nil
}
