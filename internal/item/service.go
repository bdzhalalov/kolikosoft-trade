package item

import (
	"context"
	"github.com/bdzhalalov/kolikosoft-trade/internal/item/domain"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/cache"
	customError "github.com/bdzhalalov/kolikosoft-trade/pkg/error"
	"github.com/sirupsen/logrus"
	"time"
)

type ExternalAPIClient interface {
	GetItems(ctx context.Context, params map[string]string) ([]domain.ClientResponseItem, error)
}

type Service struct {
	client ExternalAPIClient
	logger *logrus.Logger
	cache  *cache.Cache
}

func NewService(client ExternalAPIClient, logger *logrus.Logger, cache *cache.Cache) *Service {
	return &Service{
		client: client,
		logger: logger,
		cache:  cache,
	}
}

func (s *Service) GetItems(ctx context.Context) ([]GetItemsResponseDto, *customError.BaseError) {
	items, exists := s.cache.Get("items")
	if !exists {
		tradableItems, err := s.client.GetItems(ctx, nil)
		if err != nil {
			s.logger.Errorf("Error while getting tradable items: %s", err)
			return nil, (&customError.InternalServerError{}).New()
		}

		untradableItems, err := s.client.GetItems(ctx, map[string]string{
			"tradable": "0",
		})
		if err != nil {
			s.logger.Errorf("Error while getting untradable items: %s", err)
			return nil, (&customError.InternalServerError{}).New()
		}

		result := s.buildResponse(tradableItems, untradableItems)
		s.cache.Set("items", result, 5*time.Minute)

		return result, nil
	}
	return items.([]GetItemsResponseDto), nil
}

func (s *Service) buildResponse(tradable []domain.ClientResponseItem, untradable []domain.ClientResponseItem) []GetItemsResponseDto {
	items := make(map[string]*GetItemsResponseDto, len(tradable)+len(untradable))

	for i := range tradable {
		item := tradable[i]
		if item.MarketHashName == "" {
			continue
		}
		minPrice := item.MinPrice
		items[item.MarketHashName] = &GetItemsResponseDto{
			MarketHashName:     item.MarketHashName,
			Version:            item.Version,
			Currency:           item.Currency,
			SuggestedPrice:     item.SuggestedPrice,
			ItemPage:           item.ItemPage,
			MarketPage:         item.MarketPage,
			MaxPrice:           item.MaxPrice,
			MeanPrice:          item.MeanPrice,
			MedianPrice:        item.MedianPrice,
			Quantity:           item.Quantity,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
			TradableMinPrice:   &minPrice,
			UntradableMinPrice: nil,
		}
	}

	for i := range untradable {
		item := untradable[i]
		if item.MarketHashName == "" {
			continue
		}

		minPrice := item.MinPrice
		if existing, ok := items[item.MarketHashName]; ok {
			existing.UntradableMinPrice = &minPrice
			continue
		}

		items[item.MarketHashName] = &GetItemsResponseDto{
			MarketHashName:     item.MarketHashName,
			Version:            item.Version,
			Currency:           item.Currency,
			SuggestedPrice:     item.SuggestedPrice,
			ItemPage:           item.ItemPage,
			MarketPage:         item.MarketPage,
			MaxPrice:           item.MaxPrice,
			MeanPrice:          item.MeanPrice,
			MedianPrice:        item.MedianPrice,
			Quantity:           item.Quantity,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
			TradableMinPrice:   nil,
			UntradableMinPrice: &minPrice,
		}
	}

	out := make([]GetItemsResponseDto, 0, len(items))
	for _, v := range items {
		out = append(out, *v)
	}
	return out
}
