package requests

type BranchRequestDTO struct {
	Branch      string
	CountryCode string
	PhoneNumber string
	Location    string
	AddedBy     string
}

type BranchManagerRequestDTO struct {
	BranchManager string
}
