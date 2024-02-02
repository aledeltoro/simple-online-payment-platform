package models

// TransactionInput inputs to perform a transaction
type TransactionInput struct {
	Amount        int    `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Description   string `json:"description"`
}
