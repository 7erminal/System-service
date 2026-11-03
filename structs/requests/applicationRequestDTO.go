package requests

type ThemeRequest struct {
	ThemeCode string
	ThemeName string
	AddedBy   string
}

type ApplicationRequest struct {
	ApplicationName  string
	ApplicationLogo  string
	ThemeColors      string
	DefaultFontsize  string
	ApplicationImage string
	ThemeCode        string
	AddedBy          string
}

type UpdateThemeRequest struct {
	ThemeCode string
	Config    string
	UpdatedBy string
}

type ApplicationShopRequest struct {
	ShopId        string
	ApplicationId string
	AddedBy       string
}

type UpdateApplicationShopBannerRequest struct {
	ApplicationShopId int64   `json:"ApplicationShopId"`
	BannerImages      string  `json:"BannerImages"`
	BorderRadius      float64 `json:"BorderRadius"`
	UpdatedBy         string  `json:"UpdatedBy"`
}

type ToggleApplicationShopBannerRequest struct {
	ApplicationShopId int64  `json:"ApplicationShopId"`
	ShowBanner        bool   `json:"ShowBanner"`
	UpdatedBy         string `json:"UpdatedBy"`
}
