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
	Type             *string      `json:"type,omitempty"`
	UUID             *string      `json:"uuid,omitempty"`
	OrderID          *string      `json:"order_id,omitempty"`
	Amount           *string      `json:"amount,omitempty"`
	PaymentAmount    *string      `json:"payment_amount,omitempty"`
	PaymentAmountUSD *string      `json:"payment_amount_usd,omitempty"`
	MerchantAmount   *string      `json:"merchant_amount,omitempty"`
	Commission       *string      `json:"commission,omitempty"`
	IsFinal          *bool        `json:"is_final,omitempty"`
	Status           *string      `json:"status,omitempty"`
	From             *string      `json:"from,omitempty"`
	WalletAddressUUID *string     `json:"wallet_address_uuid,omitempty"`
	Network          *string      `json:"network,omitempty"`
	Currency         *string      `json:"currency,omitempty"`
	PayerCurrency    *string      `json:"payer_currency,omitempty"`
	AdditionalData   *string      `json:"additional_data,omitempty"`
	Txid             *string      `json:"txid,omitempty"`
	Convert          *ConvertData `json:"convert,omitempty"`
	Signature        *string      `json:"sign"`
}

type ConvertData struct {
	ToCurrency *string `json:"to_currency,omitempty"`
	Commission *string `json:"commission,omitempty"`
	Rate       *string `json:"rate,omitempty"`
	Amount     *string `json:"amount,omitempty"`
}
