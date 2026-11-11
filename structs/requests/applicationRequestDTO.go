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
	ApplicationShopId string  `json:"ApplicationShopId"`
	BannerImages      string  `json:"BannerImages"`
	BorderRadius      float64 `json:"BorderRadius"`
	BaseUrl           string  `json:"BaseUrl"`
	UpdatedBy         string  `json:"UpdatedBy"`
}

type ToggleApplicationShopBannerRequest struct {
	ApplicationShopId string `json:"ApplicationShopId"`
	ShowBanner        bool   `json:"ShowBanner"`
	UpdatedBy         string `json:"UpdatedBy"`
}

type UpdateApplicationThemeRequest struct {
	ApplicationId   string `json:"ApplicationId"`
	ThemeColors     string `json:"ThemeColors"`
	DefaultFontsize string `json:"DefaultFontsize"`
	UpdatedBy       string `json:"UpdatedBy"`
}
