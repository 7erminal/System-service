package responses

type Currencies struct {
	CurrencyId string
	Symbol     string
	Currency   string
}

type CurrenciesResponseDTO struct {
	StatusCode int
	Result     *[]Currencies
	StatusDesc string
}

type CurrencyResponseDTO struct {
	StatusCode int
	Result     *Currencies
	StatusDesc string
}
