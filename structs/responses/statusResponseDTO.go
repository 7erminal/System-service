package responses

type StatusResponseDTO struct {
	StatusCode string
	Status     string
	StatusId   int64
}
type StatusResponse struct {
	StatusCode int
	StatusDesc string
	Result     *StatusResponseDTO
}
type StatusListResponse struct {
	StatusCode int
	StatusDesc string
	Result     *[]StatusResponseDTO
}
