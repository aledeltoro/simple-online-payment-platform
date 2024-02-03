package models

// TransactionInput inputs to perform a transaction
type TransactionInput struct {
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	Description   string `json:"description"`
}
