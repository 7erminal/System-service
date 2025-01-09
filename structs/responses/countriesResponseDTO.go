package responses

import "system_service/models"

type CountriesResponseDTO struct {
	StatusCode int
	Countries  *[]interface{}
	StatusDesc string
}

type CountryResponseDTO struct {
	StatusCode int
	Country    *models.Countries
	StatusDesc string
}
