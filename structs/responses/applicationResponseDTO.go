package responses

import (
	"system_service/models"
	"time"
)

type ThemeResponseData struct {
	ThemeId     int64
	ThemeCode   string
	ThemeName   string
	ThemeConfig []*models.Theme_configs
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
	Theme            *ThemeResponseData
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
