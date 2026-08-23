package requests

type BillerRequestDTO struct {
	BillerName        string
	BillerCode        string
	BillerReferenceId string
	Description       string
	OperatorCode      string
	CreatedBy         string
}
