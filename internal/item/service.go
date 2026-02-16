package item

import (
	"github.com/bdzhalalov/kolikosoft-trade/internal/item/domain"
	"github.com/sirupsen/logrus"
)

type ExternalAPIClient interface {
	GetItems() ([]domain.Item, error)
}

type Service struct {
	client ExternalAPIClient
	logger *logrus.Logger
}

func NewService(client ExternalAPIClient, logger *logrus.Logger) *Service {
	return &Service{
		client: client,
		logger: logger,
	}
}
