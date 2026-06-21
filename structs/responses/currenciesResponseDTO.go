package responses

import "system_service/models"

type CurrenciesResponseDTO struct {
	StatusCode int
	Result     *[]models.Currencies
	StatusDesc string
}

type CurrencyResponseDTO struct {
	StatusCode int
	Result     *models.Currencies
	StatusDesc string
}
