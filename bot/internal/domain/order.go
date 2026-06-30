package domain

type Order struct {
	ID                string  `json:"id"`
	OrderNumber       int     `json:"order_number"`
	UserID            string  `json:"user_id"`
	RecipientUsername string  `json:"recipient_username"`
	ProductType       string  `json:"product_type"`
	ProductID         string  `json:"product_id"`
	StarsAmount       *int    `json:"stars_amount"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	Status            string  `json:"status"`
	PaymentMethod     string  `json:"payment_method"`
	PaymentURL        string  `json:"payment_url"`
}

type CatalogItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Value    int     `json:"value"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type CreateOrderInput struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	ProductType string `json:"product_type"`
	ProductID   string `json:"product_id"`
	Stars       int    `json:"stars"`
	Currency    string `json:"currency"`
}

type CreatePaymentInput struct {
	OrderID string `json:"order_id"`
	Method  string `json:"method"`
}

type WithdrawInput struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Stars    int    `json:"stars"`
}
