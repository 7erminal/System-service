package responses

import (
	"time"
)

type ThemePersonalConfigData struct {
	ThemeConfigCode string
	ThemeProperties string
	ShowBanner      bool
	BannerImages    string
	BorderRadius    float64
	DateCreated     time.Time
	DateModified    time.Time
	CreatedBy       int
	ModifiedBy      int
	Active          int
}

type ThemeConfigData struct {
	ThemeConfigCode string
	ThemeProperties string
	DateCreated     time.Time
	DateModified    time.Time
	CreatedBy       int
	ModifiedBy      int
	Active          int
}

type ThemePersonalResponseData struct {
	ThemeId      int64
	ThemeCode    string
	ThemeName    string
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
	ThemeConfig  []ThemePersonalConfigData
}

type ThemeResponseData struct {
	ThemeId      int64
	ThemeCode    string
	ThemeName    string
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
	ThemeConfig  []ThemeConfigData
}

type ThemeResponse struct {
	StatusCode    int
	StatusMessage string
	Result        *ThemeResponseData
}

type ThemeResponseList struct {
	StatusCode    int
	StatusMessage string
	Result        []ThemeResponseData
}

type ApplicationResponseData struct {
	ApplicationId    int64
	ApplicationCode  string
	ApplicationName  string
	ApplicationLogo  string
	ThemeColors      string
	DefaultFontsize  string
	ApplicationImage string
	DateCreated      time.Time
	DateModified     time.Time
	Active           int
	Theme            *ThemePersonalResponseData
	ApplicationShops []ApplicationShopResponseData
}

type ApplicationShopResponseData struct {
	ShopId string
}

type ApplicationShopFullResponseData struct {
	Application ApplicationResponseData
	ShopId      string
}

type ApplicationShopsResponse struct {
	StatusCode    int
	StatusMessage string
	Result        []ApplicationShopFullResponseData
}

type ApplicationResponse struct {
	StatusCode    int
	StatusMessage string
	Result        *ApplicationResponseData
}

type ApplicationsData struct {
	Data  *[]ApplicationResponseData
	Count int
}

type ApplicationsResponse struct {
	StatusCode    int
	StatusMessage string
	Result        ApplicationsData
}

type ApplicationThemesResponseData struct {
	Application ApplicationResponseData
	Theme       ThemeResponseData
}

type ApplicationThemesResponse struct {
	StatusCode    int
	StatusMessage string
	Result        ApplicationThemesResponseData
}

type SystemImageResponseDTO struct {
	StatusCode int
	Result     string
	StatusDesc string
}

type ShopResponse struct {
	StatusCode int
	StatusDesc string
	Result     ShopResp
}

type ShopBranchResp struct {
	ShopBranch BranchResp
	ShopId     string
	BranchId   string
}

type ShopResp struct {
	ShopId              string
	ShopName            string
	ShopDescription     string
	ShopAssistantName   string
	ShopAssistantNumber string
	PhoneNumber         string
	Email               string
	Image               string
	ShopLocation        string
	DateCreated         time.Time
	DateModified        time.Time
	CreatedBy           int
	ModifiedBy          int
	Active              int
	ShopBranches        []ShopBranchResp
}

type BranchResp struct {
	BranchId     int64
	BranchName   string
	Description  string
	Location     string
	Country      *CountryResp
	Active       int
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
}

type CountryResp struct {
	CountryId    int64
	CountryName  string
	CountryCode  string
	Active       int
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    int
	ModifiedBy   int
}
