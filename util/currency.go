package util

const (
	USD = "USD"
	EUR = "EUR"
	KZT = "KZT"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, KZT:
		return true
	}
	return false
}
