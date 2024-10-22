package services

// ExchangeI is an interface for exchanges
type ExchangeI interface {
	GetBestAdv(currency, side string, paymentMethods []string) (P2PItemI, error)
	GetName() string
	GetAds(currency, side string) ([]P2PItemI, error)
	GetAdsByName(currency, side, username string, pMethods []string) ([]P2PItemI, error)
	GetCachedPaymentMethods(curr string) ([]PaymentMethod, error)
	GetCachedCurrencies() ([]string, error)
}

// P2PItemI is an interface for exchage p2p api responses
type P2PItemI interface {
	GetPrice() float64
	GetName() string
	GetQuantity() (float64, float64, float64)
	GetPaymentMethods() []string
}

// PaymentMethod is a struct for payment methods
type PaymentMethod struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
