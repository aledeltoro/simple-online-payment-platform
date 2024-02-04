package models

// TransactionStatus type for status of transaction, defined by the payment provided
type TransactionStatus string

// PaymentProvider type for payment provider used to process transaction
type PaymentProvider string

// TransactionType type to handle type of transaction
type TransactionType string

var (
	// TransactionStatusSucceeded status for succeeded transaction
	TransactionStatusSucceeded TransactionStatus = "succeeded"
	// TransactionStatusFailure status for failed transaction
	TransactionStatusFailure TransactionStatus = "failure"
	// TransactionStatusPending status for pending transaction
	TransactionStatusPending TransactionStatus = "pending"

	// PaymentProviderStripe represents the Stripe integration
	PaymentProviderStripe PaymentProvider = "stripe"
	PaymentProviderMock   PaymentProvider = "mock"

	// TransactionTypeCharge type for processed transactions
	TransactionTypeCharge TransactionType = "charge"
	// TransactionTypeRefund type for refunded transactions
	TransactionTypeRefund TransactionType = "refund"
)

// Transaction struct to process and store a transaction
type Transaction struct {
	TransactionID    string                 `json:"transaction_id"`
	Status           TransactionStatus      `json:"status"`
	Description      string                 `json:"description"`
	FailureReason    string                 `json:"failure_reason"`
	Provider         PaymentProvider        `json:"payment_provider"`
	Amount           int                    `json:"amount"`
	Currency         string                 `json:"currency"`
	Type             TransactionType        `json:"type"`
	AdditionalFields map[string]interface{} `json:"additional_fields"`
}
