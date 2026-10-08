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
