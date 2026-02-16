package item

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/bdzhalalov/kolikosoft-trade/internal/item/domain"
	"io"
	"net/http"
	"net/url"
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

func (c *SkinPortClient) GetItems(ctx context.Context, params map[string]string) ([]domain.ClientResponseItem, error) {
	requestUrl, err := c.buildURL(fmt.Sprintf("%s/items", c.baseURL), params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Encoding", "br")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "br" {
		body = brotli.NewReader(resp.Body)
	}

	var items []domain.ClientResponseItem
	if err := json.NewDecoder(body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

func (c *SkinPortClient) buildURL(baseURL string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()

	for k, v := range params {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
