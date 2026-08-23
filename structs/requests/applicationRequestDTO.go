package requests

type ThemeRequest struct {
	ThemeCode string `json:"theme_code"`
	ThemeName string `json:"theme_name"`
}

type ApplicationRequest struct {
	ApplicationCode  string `json:"application_code"`
	ApplicationName  string `json:"application_name"`
	ApplicationLogo  string `json:"application_logo"`
	ThemeColors      string `json:"theme_colors"`
	DefaultFontsize  string `json:"default_fontsize"`
	ApplicationImage string `json:"application_image"`
	ThemeCode        string `json:"theme_code"`
}

type UpdateThemeRequest struct {
	ThemeCode string `json:"theme_code"`
	Config    string `json:"config"`
}
