package responses

import (
	"system_service/models"
)

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
