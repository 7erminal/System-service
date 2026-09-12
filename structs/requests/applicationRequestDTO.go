package requests

type ThemeRequest struct {
	ThemeCode string
	ThemeName string
}

type ApplicationRequest struct {
	ApplicationName  string
	ApplicationLogo  string
	ThemeColors      string
	DefaultFontsize  string
	ApplicationImage string
	ThemeCode        string
}

type UpdateThemeRequest struct {
	ThemeCode string
	Config    string
}
