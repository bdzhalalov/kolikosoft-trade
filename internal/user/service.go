package user

import "github.com/sirupsen/logrus"

type RepositoryInterface interface{}

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
