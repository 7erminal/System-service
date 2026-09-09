package responses

import (
	"system_service/models"
	"time"
)

type ThemeResponseData struct {
	ThemeId     int64                   `json:"theme_id"`
	ThemeCode   string                  `json:"theme_code"`
	ThemeName   string                  `json:"theme_name"`
	ThemeConfig []*models.Theme_configs `json:"theme_config"`
}

type ThemeResponse struct {
	StatusCode    int                `json:"status_code"`
	StatusMessage string             `json:"status_message"`
	Result        *ThemeResponseData `json:"result"`
}

type ThemeResponseList struct {
	StatusCode    int                 `json:"status_code"`
	StatusMessage string              `json:"status_message"`
	Result        []ThemeResponseData `json:"result"`
}

type ApplicationResponseData struct {
	ApplicationId    int64              `json:"application_id"`
	ApplicationCode  string             `json:"application_code"`
	ApplicationName  string             `json:"application_name"`
	ApplicationLogo  string             `json:"application_logo"`
	ThemeColors      string             `json:"theme_colors"`
	DefaultFontsize  string             `json:"default_fontsize"`
	ApplicationImage string             `json:"application_image"`
	DateCreated      time.Time          `json:"date_created"`
	DateModified     time.Time          `json:"date_modified"`
	Active           int                `json:"active"`
	Theme            *ThemeResponseData `json:"theme"`
}

type ApplicationResponse struct {
	StatusCode    int                      `json:"status_code"`
	StatusMessage string                   `json:"status_message"`
	Result        *ApplicationResponseData `json:"result"`
}

type ApplicationsResponse struct {
	StatusCode    int                        `json:"status_code"`
	StatusMessage string                     `json:"status_message"`
	Result        *[]ApplicationResponseData `json:"result"`
}

type ApplicationThemesResponseData struct {
	Application ApplicationResponseData `json:"application"`
	Theme       ThemeResponseData       `json:"theme"`
}

type ApplicationThemesResponse struct {
	StatusCode    int                           `json:"status_code"`
	StatusMessage string                        `json:"status_message"`
	Result        ApplicationThemesResponseData `json:"result"`
}

type SystemImageResponseDTO struct {
	StatusCode int
	Result     string
	StatusDesc string
}
