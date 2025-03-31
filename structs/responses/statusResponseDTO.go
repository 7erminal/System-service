package responses

type StatusResponseDTO struct {
	StatusCode string
	Status     string
	StatusId   int64
}
type StatusResponse struct {
	StatusCode int
	StatusDesc string
	Status     *StatusResponseDTO
}
type StatusListResponse struct {
	StatusCode int
	StatusDesc string
	Statuses   *[]StatusResponseDTO
}
