package request

type MidtransTransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

type MidtransCustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Address   string `json:"address,omitempty"`
}

type MidtransBankTransfer struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

type MidtransCreateTransactionRequest struct {
	PaymentType        string                     `json:"payment_type"`
	TransactionDetails MidtransTransactionDetails `json:"transaction_details"`
	BankTransfer       MidtransBankTransfer       `json:"bank_transfer"`
	CustomerDetails    *MidtransCustomerDetails   `json:"customer_details,omitempty"`
}

type MidtransTransactionNotificationRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"transaction_status"`
}
