package domain

type ClientResponseItem struct {
	MarketHashName string  `json:"market_hash_name"`
	Version        *string `json:"version"`
	Currency       string  `json:"currency"`
	SuggestedPrice float64 `json:"suggested_price"`
	ItemPage       string  `json:"item_page"`
	MarketPage     string  `json:"market_page"`
	MinPrice       float64 `json:"min_price"`
	MaxPrice       float64 `json:"max_price"`
	MeanPrice      float64 `json:"mean_price"`
	MedianPrice    float64 `json:"median_price"`
	Quantity       int     `json:"quantity"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
}
