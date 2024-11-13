package requests

// CryptomusRequest is used to create new invoice
type CryptomusRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	OrderID     string `json:"order_id"`
	CallbackURL string `json:"url_callback"`
	SuccessURL  string `json:"url_success"`
}

// ConfirmRequest is a webhook request from cryptomus
type ConfirmRequest struct {
	OrderID   string `json:"order_id"`
	Uuid      string `json:"uuid"`
	Status    string `json:"status"`
	Signature string `json:"sign,omitempty"`
}
