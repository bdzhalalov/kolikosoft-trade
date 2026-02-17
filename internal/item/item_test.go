package item

import (
	"context"
	"errors"
	"github.com/bdzhalalov/kolikosoft-trade/internal/item/domain"
	"github.com/bdzhalalov/kolikosoft-trade/pkg/cache"
	"github.com/jaswdr/faker/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getFakeItems() []domain.ClientResponseItem {
	f := faker.New()

	items := make([]domain.ClientResponseItem, 0, 2)

	for i := 0; i < len(items); i++ {
		items = append(items, domain.ClientResponseItem{
			MarketHashName: f.Gamer().Tag(),
			Version:        nil,
			Currency:       f.Currency().Code(),
			SuggestedPrice: f.Float64(2, 1, 50),
			ItemPage:       f.Internet().URL(),
			MarketPage:     f.Internet().URL(),
			MinPrice:       f.Float64(2, 1, 50),
			MaxPrice:       f.Float64(2, 1, 50),
			MeanPrice:      f.Float64(2, 1, 50),
			MedianPrice:    f.Float64(2, 1, 50),
			Quantity:       f.IntBetween(0, 30),
			CreatedAt:      f.Int64(),
			UpdatedAt:      f.Int64(),
		})
	}

	return items
}

type externalClientMock struct {
	err error
}

func (c *externalClientMock) GetItems(_ context.Context, params map[string]string) ([]domain.ClientResponseItem, error) {
	if c.err != nil {
		return nil, c.err
	}
	return getFakeItems(), nil
}

var (
	logger = logrus.New()
	c      = cache.New()
)

func TestGetItemsHandlerOK(t *testing.T) {
	svc := &Service{
		client: &externalClientMock{},
		logger: logger,
		cache:  c,
	}

	handler := NewHandler(svc)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/list", nil)

	handler.GetItems(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestGetItemsHandlerWithError(t *testing.T) {
	svc := &Service{
		client: &externalClientMock{err: errors.New("some error from client")},
		logger: logger,
		cache:  cache.New(),
	}

	handler := NewHandler(svc)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/list", nil)

	handler.GetItems(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestGetItemsFromCacheOK(t *testing.T) {
	svc := &Service{
		client: &externalClientMock{},
		logger: logger,
		cache:  c,
	}

	handler := NewHandler(svc)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/list", nil)

	handler.GetItems(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	_, exists := c.Get("items")
	if !exists {
		t.Fatalf("items not found in cache")
	}
}

//TODO: Add tests for the service logic (merging two lists, caching time)
