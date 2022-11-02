package request

type MidtransTransactionNotificationRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"transaction_status"`
}
