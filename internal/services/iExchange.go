package services

type ExchangeI interface {
	GetBestAdv(currency, side string, paymentMethods []string) (P2PItemI, error)
	GetName() string
	GetAdsByName(currency, side, username string) ([]P2PItemI, error)
}

type P2PItemI interface {
	GetPrice() float64
	GetName() string
	GetQuantity() (float64, float64, float64)
	GetPaymentMethods() []string
}
