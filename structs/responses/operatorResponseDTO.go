package responses

type OperatorResponseDTO struct {
	StatusCode int
	Operator   *OperatorObject
	StatusDesc string
}

type OperatorObject struct {
	OperatorId   string
	OperatorName string
	Description  string
	DateCreated  string
	DateModified string
	CreatedBy    string
	ModifiedBy   string
	Active       int
}

type OperatorsResponseDTO struct {
	StatusCode int
	Operator   []OperatorObject
	StatusDesc string
}
