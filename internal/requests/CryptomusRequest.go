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
	Type             *string      `json:"type"`
	UUID             *string      `json:"uuid"`
	OrderID          *string      `json:"order_id"`
	Amount           *string      `json:"amount"`
	PaymentAmount    *string      `json:"payment_amount"`
	PaymentAmountUSD *string      `json:"payment_amount_usd"`
	MerchantAmount   *string      `json:"merchant_amount"`
	Commission       *string      `json:"commission"`
	IsFinal          *bool        `json:"is_final"`
	Status           *string      `json:"status"`
	From             *string      `json:"from"`
	WalletAddressUUID *string     `json:"wallet_address_uuid"`
	Network          *string      `json:"network"`
	Currency         *string      `json:"currency"`
	PayerCurrency    *string      `json:"payer_currency"`
	AdditionalData   *string      `json:"additional_data"`
	Txid             *string      `json:"txid"`
	Convert          *ConvertData `json:"convert"`
	Signature        *string      `json:"sign"`
}

type ConvertData struct {
	ToCurrency *string `json:"to_currency"`
	Commission *string `json:"commission"`
	Rate       *string `json:"rate"`
	Amount     *string `json:"amount"`
}
