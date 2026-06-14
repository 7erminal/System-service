package responses

import (
	"system_service/models"
	"time"
)

type Branches struct {
	BranchId      int64
	Branch        string
	Country       *models.Countries
	Location      string
	PhoneNumber   string
	Active        int
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     int
	ModifiedBy    int
	BranchManager *Users
}

type BranchesResponseDTO struct {
	StatusCode int
	Branches   *[]interface{}
	StatusDesc string
}

type BranchResponseDTO struct {
	StatusCode int
	Branch     *Branches
	StatusDesc string
}
