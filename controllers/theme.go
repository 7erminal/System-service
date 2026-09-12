package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"system_service/models"
	"system_service/structs/requests"
	"system_service/structs/responses"
	"time"

	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
)

// ThemeController operations for Theme
type ThemeController struct {
	beego.Controller
}

// URLMapping ...
func (c *ThemeController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("AddThemeConfig", c.AddThemeConfig)
	c.Mapping("RemoveThemeConfig", c.RemoveThemeConfig)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Theme
// @Param	body		body 	requests.ThemeRequest	true		"body for Theme content"
// @Success 201 {int} responses.ThemeResponse
// @Failure 403 body is empty
// @router / [post]
func (c *ThemeController) Post() {
	var req requests.ThemeRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ThemeResponseData

	logs.Info("Request received: ", req)

	v := models.Theme{ThemeCode: req.ThemeCode, ThemeName: req.ThemeName}
	if _, err := models.AddTheme(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		statusCode = 200
		statusMessage = "Theme created successfully. Creating configuration for the theme..."

		themeConfig := models.Theme_configs{
			ThemeConfigCode: req.ThemeCode + "_CONFIG",
			ThemeProperties: "",
			ThemeId:         &v,
		}

		if _, err := models.AddTheme_configs(&themeConfig); err == nil {
			statusMessage += " Theme configuration created successfully."
		} else {
			statusCode = 500
			statusMessage = "Internal server error: " + err.Error()
		}
		result = responses.ThemeResponseData{
			ThemeId:   v.ThemeId,
			ThemeCode: v.ThemeCode,
			ThemeName: v.ThemeName,
			ThemeConfig: []*models.Theme_configs{
				&themeConfig,
			},
		}

	} else {
		c.Data["json"] = err.Error()
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
	}
	response := responses.ThemeResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &result,
	}
	logs.Info("Response in json:: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Theme by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} responses.ThemeResponse
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ThemeController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	logs.Info("Request received: id=", idStr)
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetThemeById(id)
	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ThemeResponseData

	if err == nil {
		statusCode = 200
		statusMessage = "Theme retrieved successfully"
		result = responses.ThemeResponseData{
			ThemeId:     v.ThemeId,
			ThemeCode:   v.ThemeCode,
			ThemeName:   v.ThemeName,
			ThemeConfig: v.ThemeConfigs,
		}
	} else {
		statusCode = 404
		statusMessage = "Theme not found"
	}

	response := responses.ThemeResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &result,
	}
	logs.Info("Response in json:: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Theme
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} responses.ThemeResponseList
// @Failure 403
// @router / [get]
func (c *ThemeController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	logs.Info("GetAll request received")

	statusCode := 400
	statusMessage := "Bad Request"
	var result []responses.ThemeResponseData

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllTheme(query, fields, sortby, order, offset, limit)
	if err != nil {
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
	} else {
		statusCode = 200
		statusMessage = "Themes retrieved successfully"

		themesResp := []responses.ThemeResponseData{}
		for _, urs := range l {
			m := urs.(models.Theme)

			themesResp = append(themesResp, responses.ThemeResponseData{
				ThemeId:     m.ThemeId,
				ThemeCode:   m.ThemeCode,
				ThemeName:   m.ThemeName,
				ThemeConfig: m.ThemeConfigs,
			})
		}

		result = themesResp
	}

	response := responses.ThemeResponseList{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        result,
	}
	logs.Info("GetAll response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Theme
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	requests.ThemeRequest	true		"body for Theme content"
// @Success 200 {object} models.Theme
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ThemeController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var req requests.UpdateThemeRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	logs.Info("Put request received: id=", idStr, " body=", req)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ThemeResponseData

	// Check if the theme exists
	_, err := models.GetThemeById(id)
	if err != nil {
		statusCode = 404
		statusMessage = "Theme not found"
		response := responses.ThemeResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		logs.Info("Put response: ", response)
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	theme, err := models.GetThemeById(id)
	if err != nil {
		statusCode = 404
		statusMessage = "Theme not found"
	} else {
		themeConfig, err := models.GetTheme_configsByThemeId(theme.ThemeId)
		if err != nil {
			statusCode = 404
			statusMessage = "Theme configuration not found"
		} else {
			themeConfig.ThemeProperties = req.Config
			if err := models.UpdateTheme_configsById(themeConfig); err == nil {
				statusCode = 200
				statusMessage = "Theme updated successfully"
				result = responses.ThemeResponseData{
					ThemeId:   theme.ThemeId,
					ThemeCode: theme.ThemeCode,
					ThemeName: theme.ThemeName,
					ThemeConfig: []*models.Theme_configs{
						themeConfig,
					},
				}
			} else {
				statusCode = 500
				statusMessage = "Internal server error: " + err.Error()
			}
		}
	}

	response := responses.ThemeResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &result,
	}
	logs.Info("Put response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// AddThemeConfig ...
// @Title Add Theme Config
// @Description add a theme configuration for an existing theme
// @Param	id		path 	string	true		"The theme id"
// @Param	body		body 	requests.UpdateThemeRequest	true		"body for theme config"
// @Success 200 {object} responses.ThemeResponse
// @Failure 400,404,500
// @router /:id/config [post]
func (c *ThemeController) AddThemeConfig() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	var req requests.UpdateThemeRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ThemeResponseData

	if err != nil {
		statusMessage = "Invalid theme id"
		response := responses.ThemeResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &result}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(req.Config) == "" {
		statusMessage = "config is required"
		response := responses.ThemeResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &result}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	theme, err := models.GetThemeById(id)
	if err != nil {
		statusCode = 404
		statusMessage = "Theme not found"
		response := responses.ThemeResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &result}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	themeConfig := models.Theme_configs{
		ThemeId:         theme,
		ThemeConfigCode: theme.ThemeCode + "_CONFIG_" + strconv.FormatInt(time.Now().UnixNano(), 10),
		ThemeProperties: req.Config,
	}

	if _, err := models.AddTheme_configs(&themeConfig); err != nil {
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
	} else {
		freshTheme, getErr := models.GetThemeById(id)
		if getErr != nil {
			statusCode = 500
			statusMessage = "Theme config added but failed to fetch updated theme"
		} else {
			statusCode = 200
			statusMessage = "Theme config added successfully"
			result = responses.ThemeResponseData{
				ThemeId:     freshTheme.ThemeId,
				ThemeCode:   freshTheme.ThemeCode,
				ThemeName:   freshTheme.ThemeName,
				ThemeConfig: freshTheme.ThemeConfigs,
			}
		}
	}

	response := responses.ThemeResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &result}
	logs.Info("AddThemeConfig response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// RemoveThemeConfig ...
// @Title Remove Theme Config
// @Description remove a theme configuration by config id
// @Param	id		path 	string	true		"The theme config id"
// @Success 200 {object} responses.ThemeResponse
// @Failure 400,404,500
// @router /config/:id [delete]
func (c *ThemeController) RemoveThemeConfig() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ThemeResponseData

	if err != nil {
		statusMessage = "Invalid theme config id"
		response := responses.ThemeResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &result}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	themeConfig, err := models.GetTheme_configsById(id)
	if err != nil {
		statusCode = 404
		statusMessage = "Theme configuration not found"
		response := responses.ThemeResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &result}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	theme, err := models.GetThemeById(themeConfig.ThemeId.ThemeId)
	if err != nil {
		statusCode = 404
		statusMessage = "Theme not found"
		response := responses.ThemeResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &result}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	if err := models.DeleteTheme_configs(id); err != nil {
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
	} else {
		freshTheme, getErr := models.GetThemeById(theme.ThemeId)
		if getErr != nil {
			statusCode = 500
			statusMessage = "Theme config removed but failed to fetch updated theme"
		} else {
			statusCode = 200
			statusMessage = "Theme config removed successfully"
			result = responses.ThemeResponseData{
				ThemeId:     freshTheme.ThemeId,
				ThemeCode:   freshTheme.ThemeCode,
				ThemeName:   freshTheme.ThemeName,
				ThemeConfig: freshTheme.ThemeConfigs,
			}
		}
	}

	response := responses.ThemeResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &result}
	logs.Info("RemoveThemeConfig response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Theme
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ThemeController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	logs.Info("Delete request received: id=", idStr)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ThemeResponseData

	// Check if the theme exists
	_, err := models.GetThemeById(id)
	if err != nil {
		statusCode = 404
		statusMessage = "Theme not found"
		response := responses.ThemeResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		logs.Info("Delete response: ", response)
		c.Data["json"] = response
		c.ServeJSON()
		return
	}
	if err := models.DeleteTheme(id); err == nil {
		statusCode = 200
		statusMessage = "Theme deleted successfully"
		result = responses.ThemeResponseData{
			ThemeId: id,
		}

	} else {
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
	}
	response := responses.ThemeResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &result,
	}
	logs.Info("Delete response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}
