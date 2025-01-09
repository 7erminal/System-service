package responses

import "system_service/models"

type CurrenciesResponseDTO struct {
	StatusCode int
	Currencies *[]interface{}
	StatusDesc string
}

type CurrencyResponseDTO struct {
	StatusCode int
	Currency   *models.Currencies
	StatusDesc string
}
