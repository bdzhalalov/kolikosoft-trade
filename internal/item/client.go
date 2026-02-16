package item

import (
	"github.com/bdzhalalov/kolikosoft-trade/internal/item/domain"
	"net/http"
)

type SkinPortClient struct {
	client  *http.Client
	baseURL string
}

func NewSkinPortClient(client *http.Client, baseURL string) *SkinPortClient {
	return &SkinPortClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *SkinPortClient) GetItems() ([]domain.Item, error) {
	var items []domain.Item
	return items, nil
}
