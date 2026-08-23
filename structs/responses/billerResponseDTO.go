package responses

type BillerResponseDTO struct {
	StatusCode int
	Biller     *BillerObject
	StatusDesc string
}

type BillerObject struct {
	BillerId          int64
	BillerName        string
	BillerCode        string
	BillerReferenceId string
	Description       string
	Operator          string
	DateCreated       string
	DateModified      string
	CreatedBy         string
	ModifiedBy        string
	Active            int
}

type BillersResponseDTO struct {
	StatusCode int
	Result     []BillerObject
	StatusDesc string
}
