package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"system_service/models"
	"system_service/structs/requests"
	"system_service/structs/responses"

	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
)

// ApplicationController operations for Application
type ApplicationController struct {
	beego.Controller
}

// URLMapping ...
func (c *ApplicationController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Application
// @Param	body		body 	requests.ApplicationRequest	true		"body for Application content"
// @Success 201 {int} responses.ApplicationResponse
// @Failure 403 body is empty
// @router / [post]
func (c *ApplicationController) Post() {
	var req requests.ApplicationRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ApplicationResponseData

	logs.Info("Post request received: ", req)

	theme, err := models.GetThemeByCode(req.ThemeCode)
	if err != nil {
		statusCode = 404
		statusMessage = "Theme not found"
	} else {

		v := models.Application{
			ApplicationCode:  req.ApplicationCode,
			ApplicationName:  req.ApplicationName,
			ApplicationLogo:  req.ApplicationLogo,
			ThemeColors:      req.ThemeColors,
			DefaultFontsize:  req.DefaultFontsize,
			ApplicationImage: req.ApplicationImage,
		}
		if id, err := models.AddApplication(&v); err == nil {

			statusCode = 200
			statusMessage = "Application created. Adding theme association."

			apt := models.Application_themes{
				ApplicationId: &v,
				ThemeId:       theme,
			}
			if apthid, err := models.AddApplication_themes(&apt); err == nil {
				logs.Info("ApplicationTheme created with ID: ", apthid)

				c.Ctx.Output.SetStatus(201)
				statusCode = 200
				statusMessage = "Application created successfully"
				v.ApplicationId = id

				themeResp := responses.ThemeResponseData{
					ThemeId:   theme.ThemeId,
					ThemeCode: theme.ThemeCode,
					ThemeName: theme.ThemeName,
				}
				result = responses.ApplicationResponseData{
					ApplicationId:    v.ApplicationId,
					ApplicationCode:  v.ApplicationCode,
					ApplicationName:  v.ApplicationName,
					ApplicationLogo:  v.ApplicationLogo,
					ThemeColors:      v.ThemeColors,
					DefaultFontsize:  v.DefaultFontsize,
					ApplicationImage: v.ApplicationImage,
					DateCreated:      v.DateCreated,
					DateModified:     v.DateModified,
					Active:           v.Active,
					Theme:            &themeResp,
				}
			} else {
				logs.Error("Error creating ApplicationTheme: ", err)
			}
		} else {
			statusCode = 500
			statusMessage = "Internal server error: " + err.Error()
		}
	}

	response := responses.ApplicationResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &result,
	}
	logs.Info("Post response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Application by id
// @Param	id		path 	string	true		"The key for Application"
// @Success 200 {object} responses.ApplicationResponse
// @Failure 403 :id is empty
// @router /:code [get]
func (c *ApplicationController) GetOne() {
	code := c.Ctx.Input.Param(":code")
	logs.Info("GetOne request received: code=", code)
	a, err := models.GetApplicationByCode(code)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ApplicationResponseData

	if err != nil {

	} else {

		v, erra := models.GetApplication_themesByApplicationId(a.ApplicationId)
		if erra == nil {
			statusCode = 200
			statusMessage = "Application retrieved successfully"
			result = responses.ApplicationResponseData{
				ApplicationId:    v.ApplicationId.ApplicationId,
				ApplicationCode:  v.ApplicationId.ApplicationCode,
				ApplicationName:  v.ApplicationId.ApplicationName,
				ApplicationLogo:  v.ApplicationId.ApplicationLogo,
				ThemeColors:      v.ApplicationId.ThemeColors,
				DefaultFontsize:  v.ApplicationId.DefaultFontsize,
				ApplicationImage: v.ApplicationId.ApplicationImage,
				DateCreated:      v.ApplicationId.DateCreated,
				DateModified:     v.ApplicationId.DateModified,
				Active:           v.ApplicationId.Active,
				Theme: &responses.ThemeResponseData{
					ThemeId:   v.ThemeId.ThemeId,
					ThemeCode: v.ThemeId.ThemeCode,
					ThemeName: v.ThemeId.ThemeName,
				},
			}
		} else {
			statusCode = 404
			statusMessage = "Application not found"
		}
	}

	response := responses.ApplicationResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &result,
	}
	logs.Info("GetOne response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Application
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} responses.ApplicationsResponse
// @Failure 403
// @router / [get]
func (c *ApplicationController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	logs.Info("GetAll request received")

	statusCode := 400
	statusMessage := "Bad Request"
	var result []responses.ApplicationResponseData

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

	l, err := models.GetAllApplication_themes(query, fields, sortby, order, offset, limit)
	if err != nil {
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
	} else {
		statusCode = 200
		statusMessage = "Applications retrieved successfully"

		appsResp := []responses.ApplicationResponseData{}
		for _, app := range l {
			m := app.(models.Application_themes)
			appsResp = append(appsResp, responses.ApplicationResponseData{
				ApplicationId:    m.ApplicationId.ApplicationId,
				ApplicationCode:  m.ApplicationId.ApplicationCode,
				ApplicationName:  m.ApplicationId.ApplicationName,
				ApplicationLogo:  m.ApplicationId.ApplicationLogo,
				ThemeColors:      m.ApplicationId.ThemeColors,
				DefaultFontsize:  m.ApplicationId.DefaultFontsize,
				ApplicationImage: m.ApplicationId.ApplicationImage,
				DateCreated:      m.ApplicationId.DateCreated,
				DateModified:     m.ApplicationId.DateModified,
				Active:           m.ApplicationId.Active,
				Theme: &responses.ThemeResponseData{
					ThemeId:   m.ThemeId.ThemeId,
					ThemeCode: m.ThemeId.ThemeCode,
					ThemeName: m.ThemeId.ThemeName,
				},
			})
		}
		result = appsResp
	}

	responseList := result
	response := responses.ApplicationsResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &responseList,
	}
	logs.Info("GetAll response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Application
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	requests.ApplicationRequest	true		"body for Application content"
// @Success 200 {object} responses.ApplicationResponse
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ApplicationController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var req requests.ApplicationRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	logs.Info("Put request received: id=", idStr, " body=", req)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ApplicationResponseData

	// Check if the application exists
	_, err := models.GetApplicationById(id)
	if err != nil {
		statusCode = 404
		statusMessage = "Application not found"
		response := responses.ApplicationResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		logs.Info("Put response: ", response)
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	v := models.Application{
		ApplicationId:    id,
		ApplicationCode:  req.ApplicationCode,
		ApplicationName:  req.ApplicationName,
		ApplicationLogo:  req.ApplicationLogo,
		ThemeColors:      req.ThemeColors,
		DefaultFontsize:  req.DefaultFontsize,
		ApplicationImage: req.ApplicationImage,
	}
	if err := models.UpdateApplicationById(&v); err == nil {
		statusCode = 200
		statusMessage = "Application updated successfully"
		result = responses.ApplicationResponseData{
			ApplicationId:    v.ApplicationId,
			ApplicationCode:  v.ApplicationCode,
			ApplicationName:  v.ApplicationName,
			ApplicationLogo:  v.ApplicationLogo,
			ThemeColors:      v.ThemeColors,
			DefaultFontsize:  v.DefaultFontsize,
			ApplicationImage: v.ApplicationImage,
			DateCreated:      v.DateCreated,
			DateModified:     v.DateModified,
			Active:           v.Active,
		}
	} else {
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
	}

	response := responses.ApplicationResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &result,
	}
	logs.Info("Put response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Application
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {object} responses.ApplicationResponse
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ApplicationController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	logs.Info("Delete request received: id=", idStr)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ApplicationResponseData

	// Check if the application exists
	_, err := models.GetApplicationById(id)
	if err != nil {
		statusCode = 404
		statusMessage = "Application not found"
		response := responses.ApplicationResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		logs.Info("Delete response: ", response)
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	if err := models.DeleteApplication(id); err == nil {
		statusCode = 200
		statusMessage = "Application deleted successfully"
		result = responses.ApplicationResponseData{
			ApplicationId: id,
		}
	} else {
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
	}

	response := responses.ApplicationResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &result,
	}
	logs.Info("Delete response: ", response)
	c.Data["json"] = response
	c.ServeJSON()
}
