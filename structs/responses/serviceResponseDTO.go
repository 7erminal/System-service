package responses

type ServiceResponseDTO struct {
	StatusCode int
	Result     ServiceObject
	StatusDesc string
}

type ServiceObject struct {
	ServiceId          int64
	ServiceName        string
	ServiceCode        string
	ServiceDescription string
	DateCreated        string
	DateModified       string
	CreatedBy          int
	ModifiedBy         int
	Active             int
}

type ServicesResponseDTO struct {
	StatusCode int
	Result     []ServiceObject
	StatusDesc string
}
