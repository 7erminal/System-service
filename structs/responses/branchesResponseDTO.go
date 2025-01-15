package responses

import (
	"system_service/models"
	"time"
)

type BranchesResp struct {
	BranchId      int64             `orm:"auto"`
	Branch        string            `orm:"size(80)"`
	Country       *models.Countries `orm:"rel(fk);column(country_id)"`
	Location      string            `orm:"column(location)"`
	PhoneNumber   string            `orm:"column(phone_number)"`
	Active        int               `orm:"omitempty"`
	DateCreated   time.Time         `orm:"type(datetime);omitempty"`
	DateModified  time.Time         `orm:"type(datetime);omitempty"`
	CreatedBy     int               `orm:"omitempty"`
	ModifiedBy    int               `orm:"omitempty"`
	BranchManager *models.User
}

type BranchesResponseDTO struct {
	StatusCode int
	Branches   *[]interface{}
	StatusDesc string
}

type BranchResponseDTO struct {
	StatusCode int
	Branch     *models.Branches
	StatusDesc string
}

type BranchResponseRDTO struct {
	StatusCode int
	Branch     *BranchesResp
	StatusDesc string
}
