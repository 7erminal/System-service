package responses

import "system_service/models"

type CountriesResponseDTO struct {
	StatusCode int
	Result     *[]models.Countries
	StatusDesc string
}

type CountryResponseDTO struct {
	StatusCode int
	Result     *models.Countries
	StatusDesc string
}
