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
	AdditionalData           *string     `json:"additional_data"`
	Amount                   *string     `json:"amount"`
	Commission               *string     `json:"commission"`
	Currency                 *string     `json:"currency"`
	From                     *string     `json:"from"`
	IsFinal                  *bool       `json:"is_final"`
	MerchantAmount           *string     `json:"merchant_amount"`
	Network                  *string     `json:"network"`
	OrderID                  *string     `json:"order_id"`
	PayerAmount              *string     `json:"payer_amount"`
	PayerAmountExchangeRate  *string     `json:"payer_amount_exchange_rate"`
	PayerCurrency            *string     `json:"payer_currency"`
	PaymentAmount            *string     `json:"payment_amount"`
	PaymentAmountUSD         *string     `json:"payment_amount_usd"`
	Signature                *string     `json:"sign,omitempty"`
	Status                   *string     `json:"status"`
	TransferID               *string     `json:"transfer_id"`
	Txid                     *string     `json:"txid"`
	Type                     *string     `json:"type"`
	UUID                     *string     `json:"uuid"`
	WalletAddressUUID        *string     `json:"wallet_address_uuid"`
	Convert                  *ConvertData `json:"convert,omitempty"`  // Optional nested struct for `convert`
}

type ConvertData struct {
	ToCurrency *string `json:"to_currency"`
	Commission *string `json:"commission"`
	Rate       *string `json:"rate"`
	Amount     *string `json:"amount"`
}
