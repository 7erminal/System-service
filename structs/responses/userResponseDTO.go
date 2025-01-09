package responses

type UserResponseDTO struct {
	StatusCode int
	User       *interface{}
	StatusDesc string
}

type UserGatewayResponseDTO struct {
	Success    bool
	Result     *interface{}
	StatusDesc string
}
