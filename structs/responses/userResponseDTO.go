package responses

import (
	"time"
)

type Shops struct {
	ShopId              int64
	ShopName            string
	ShopDescription     string
	ShopAssistantName   string
	ShopAssistantNumber string
	PhoneNumber         string
	Email               string
	Image               string
	DateCreated         time.Time
	DateModified        time.Time
	CreatedBy           int
	ModifiedBy          int
	Active              int
}

type UserExtraDetails struct {
	UserDetailsId int64
	Branch        *Branches
	Shop          *Shops
	Nickname      string
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     int
	ModifiedBy    int
	Active        int
}

type Users struct {
	UserId        int64
	UserDetails   *UserExtraDetails
	ImagePath     string
	UserType      int
	FullName      string
	Username      string
	Password      string
	Email         string
	PhoneNumber   string
	Gender        string
	Dob           time.Time
	Address       string
	IdType        string
	IdNumber      string
	MaritalStatus string
	Active        int
	Role          *Roles
	IsVerified    bool
	DateCreated   time.Time
	DateModified  time.Time
	CreatedBy     int
	ModifiedBy    int
}

type Roles struct {
	RoleId       int64
	Role         string
	Description  string
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
	Active       int
}

type UserResponseDTO struct {
	StatusCode int
	User       *Users
	StatusDesc string
}

type UserGatewayResponseDTO struct {
	Success    bool
	Result     *interface{}
	StatusDesc string
}
