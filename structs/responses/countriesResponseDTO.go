package responses

import (
	"time"
)

type Countries struct {
	CountryId       string
	Country         string
	Description     string
	CountryCode     string
	DefaultCurrency *Currencies
	DateCreated     time.Time
	DateModified    time.Time
	CreatedBy       string
	ModifiedBy      string
}

type CountriesResponseDTO struct {
	StatusCode int
	Result     *[]Countries
	StatusDesc string
}

type CountryResponseDTO struct {
	StatusCode int
	Result     *Countries
	StatusDesc string
}
