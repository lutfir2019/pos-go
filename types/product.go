package types

type Products struct {
	ID            uint   `json:"id"`
	Code          string `json:"product_code"`
	Name          string `json:"name"`
	Quantity      uint   `json:"quantity"`
	PricePurchase uint   `json:"price_purchase"`
	PriceSelling  uint   `json:"price_selling"`
	Image         string `json:"file"`
	CreatedAt     string `json:"created_at"`
	UpdateedAt    string `json:"updated_at"`
	ReferUser     string `json:"created_by"`
}
