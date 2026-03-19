package responses

import "system_service/models"

type UserResponseDTO struct {
	StatusCode int
	User       *models.Users
	StatusDesc string
}

type UserGatewayResponseDTO struct {
	Success    bool
	Result     *interface{}
	StatusDesc string
}
