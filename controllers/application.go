package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"system_service/controllers/functions"
	"system_service/models"
	"system_service/structs/requests"
	"system_service/structs/responses"
	"time"

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
	c.Mapping("UpdateTheme", c.UpdateTheme)
	c.Mapping("Delete", c.Delete)
	c.Mapping("UploadImage", c.UploadImage)
	c.Mapping("RemoveApplicationShop", c.RemoveApplicationShop)
	c.Mapping("AddApplicationShop", c.AddApplicationShop)
	c.Mapping("GetApplicationShops", c.GetApplicationShops)
	c.Mapping("UpdateApplicationShopBanner", c.UpdateApplicationShopBanner)
	c.Mapping("ToggleApplicationShopBanner", c.ToggleApplicationShopBanner)
	c.Mapping("UpdateApplicationTheme", c.UpdateApplicationTheme)
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

		// Generate application code
		applicationCode := fmt.Sprintf("APP-%d", time.Now().Unix())
		v := models.Application{
			ApplicationCode:  applicationCode,
			ApplicationName:  req.ApplicationName,
			ApplicationLogo:  req.ApplicationLogo,
			ThemeColors:      req.ThemeColors,
			DefaultFontsize:  req.DefaultFontsize,
			ApplicationImage: req.ApplicationImage,
			Active:           1,
		}
		if id, err := models.AddApplication(&v); err == nil {

			statusCode = 200
			statusMessage = "Application created. Adding theme association."

			apt := models.Application_themes{
				Application: &v,
				Theme:       theme,
			}
			if apthid, err := models.AddApplication_themes(&apt); err == nil {
				logs.Info("ApplicationTheme created with ID: ", apthid)

				c.Ctx.Output.SetStatus(201)
				statusCode = 200
				statusMessage = "Application created successfully"
				v.ApplicationId = id

				themeResp := responses.ThemePersonalResponseData{
					ThemeId:   strconv.FormatInt(theme.ThemeId, 10),
					ThemeCode: theme.ThemeCode,
					ThemeName: theme.ThemeName,
				}
				result = responses.ApplicationResponseData{
					ApplicationId:    strconv.FormatInt(v.ApplicationId, 10),
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
		statusCode = 404
		statusMessage = "Application not found"
	} else {
		logs.Info("Application retrieved: ", a)
		logs.Info("Application shops retrieved: ", a.ApplicationShops)
		v, erra := models.GetApplication_themesByApplicationId(a.ApplicationId)
		if erra == nil {
			statusCode = 200
			statusMessage = "Application retrieved successfully"

			themeConfigs := []responses.ThemePersonalConfigData{}

			for _, tc := range v.Theme.ThemeConfigs {
				themeConfigs = append(themeConfigs, responses.ThemePersonalConfigData{
					ThemeConfigCode: tc.ThemeConfigCode,
					ThemeProperties: tc.ThemeProperties,
					DateCreated:     tc.DateCreated,
					DateModified:    tc.DateModified,
					CreatedBy:       tc.CreatedBy,
					ModifiedBy:      tc.ModifiedBy,
					Active:          tc.Active,
				})
			}
			themeResp := responses.ThemePersonalResponseData{
				ThemeId:      strconv.FormatInt(v.Theme.ThemeId, 10),
				ThemeCode:    v.Theme.ThemeCode,
				ThemeName:    v.Theme.ThemeName,
				ThemeConfig:  themeConfigs,
				DateCreated:  v.Theme.DateCreated,
				DateModified: v.Theme.DateModified,
				CreatedBy:    v.Theme.CreatedBy,
				ModifiedBy:   v.Theme.ModifiedBy,
			}

			logs.Info("Application shops retrieved :: ", a.ApplicationShops)
			appShops := []responses.ApplicationShopResponseData{}
			for _, as := range a.ApplicationShops {
				appShops = append(appShops, responses.ApplicationShopResponseData{
					ShopId:            as.Shop,
					ApplicationId:     strconv.FormatInt(a.ApplicationId, 10),
					ApplicationShopId: strconv.FormatInt(as.Id, 10),
					ApplicationCode:   a.ApplicationCode,
					ApplicationName:   a.ApplicationName,
					ApplicationLogo:   a.ApplicationLogo,
					ApplicationImage:  a.ApplicationImage,
					ShowBanner:        as.ShowBanner,
					BannerImages:      as.BannerImages,
					BorderRadius:      strconv.FormatFloat(as.BorderRadius, 'f', -1, 64),
					ThemeColors:       a.ThemeColors,
					DefaultFontsize:   a.DefaultFontsize,
					DateCreated:       a.DateCreated,
					DateModified:      a.DateModified,
					Active:            a.Active,
				})
			}
			result = responses.ApplicationResponseData{
				ApplicationId:    strconv.FormatInt(a.ApplicationId, 10),
				ApplicationCode:  a.ApplicationCode,
				ApplicationName:  a.ApplicationName,
				ApplicationLogo:  a.ApplicationLogo,
				ThemeColors:      a.ThemeColors,
				DefaultFontsize:  a.DefaultFontsize,
				ApplicationImage: a.ApplicationImage,
				DateCreated:      a.DateCreated,
				DateModified:     a.DateModified,
				Active:           a.Active,
				Theme:            &themeResp,
				Shops:            appShops,
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

			themeConfigs := []responses.ThemePersonalConfigData{}

			for _, tc := range m.Theme.ThemeConfigs {
				themeConfigs = append(themeConfigs, responses.ThemePersonalConfigData{
					ThemeConfigCode: tc.ThemeConfigCode,
					ThemeProperties: tc.ThemeProperties,
					DateCreated:     tc.DateCreated,
					DateModified:    tc.DateModified,
					CreatedBy:       tc.CreatedBy,
					ModifiedBy:      tc.ModifiedBy,
					Active:          tc.Active,
				})
			}

			appShops := []responses.ApplicationShopResponseData{}
			shopQuery := make(map[string]string)
			shopQuery["Application"] = strconv.FormatInt(m.Application.ApplicationId, 10)
			if shops, err := models.GetAllApplicationShops(shopQuery, []string{}, []string{}, []string{}, 0, 1000); err == nil {
				for _, shop := range shops {
					s := shop.(models.ApplicationShops)
					appShops = append(appShops, responses.ApplicationShopResponseData{
						ApplicationShopId: strconv.FormatInt(s.Id, 10),
						ApplicationId:     strconv.FormatInt(m.Application.ApplicationId, 10),
						ApplicationCode:   m.Application.ApplicationCode,
						ApplicationName:   m.Application.ApplicationName,
						ApplicationLogo:   m.Application.ApplicationLogo,
						ShopId:            s.Shop,
						ShowBanner:        s.ShowBanner,
						BannerImages:      s.BannerImages,
						BorderRadius:      strconv.FormatFloat(s.BorderRadius, 'f', -1, 64),
						ThemeColors:       m.Application.ThemeColors,
						DefaultFontsize:   m.Application.DefaultFontsize,
						ApplicationImage:  m.Application.ApplicationImage,
						DateCreated:       s.DateCreated,
						DateModified:      s.DateModified,
						CreatedBy:         s.CreatedBy,
						ModifiedBy:        s.ModifiedBy,
						Active:            s.Active,
					})
				}
			}

			appsResp = append(appsResp, responses.ApplicationResponseData{
				ApplicationId:    strconv.FormatInt(m.Application.ApplicationId, 10),
				ApplicationCode:  m.Application.ApplicationCode,
				ApplicationName:  m.Application.ApplicationName,
				ApplicationLogo:  m.Application.ApplicationLogo,
				ThemeColors:      m.Application.ThemeColors,
				DefaultFontsize:  m.Application.DefaultFontsize,
				ApplicationImage: m.Application.ApplicationImage,
				DateCreated:      m.Application.DateCreated,
				DateModified:     m.Application.DateModified,
				Active:           m.Application.Active,
				Theme: &responses.ThemePersonalResponseData{
					ThemeId:      strconv.FormatInt(m.Theme.ThemeId, 10),
					ThemeCode:    m.Theme.ThemeCode,
					ThemeName:    m.Theme.ThemeName,
					DateCreated:  m.Theme.DateCreated,
					DateModified: m.Theme.DateModified,
					CreatedBy:    m.Theme.CreatedBy,
					ModifiedBy:   m.Theme.ModifiedBy,
					ThemeConfig:  themeConfigs,
				},
				Shops: appShops,
			})
		}
		result = appsResp
	}

	responseList := result
	logs.Info("Number of applications retrieved: ", len(responseList))

	respData := responses.ApplicationsData{
		Data:  &responseList,
		Count: len(responseList),
	}
	response := responses.ApplicationsResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        respData,
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
	app, err := models.GetApplicationById(id)
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
		ApplicationCode:  app.ApplicationCode,
		ApplicationName:  req.ApplicationName,
		ApplicationLogo:  req.ApplicationLogo,
		ThemeColors:      req.ThemeColors,
		DefaultFontsize:  req.DefaultFontsize,
		ApplicationImage: req.ApplicationImage,
		Active:           app.Active,
		ModifiedBy:       req.AddedBy,
	}
	if err := models.UpdateApplicationById(&v); err == nil {
		statusCode = 200
		statusMessage = "Application updated successfully"
		result = responses.ApplicationResponseData{
			ApplicationId:    strconv.FormatInt(v.ApplicationId, 10),
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

// UpdateTheme ...
// @Title Update Application Theme
// @Description update application theme by application id using theme code
// @Param	id		path 	string	true		"The application id"
// @Param	body		body 	requests.UpdateThemeRequest	true		"body with theme_code"
// @Success 200 {object} responses.ApplicationResponse
// @Failure 400,404,500
// @router /:id/theme [put]
func (c *ApplicationController) UpdateTheme() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)

	statusCode := 400
	statusMessage := "Bad Request"
	var result responses.ApplicationResponseData

	if err != nil {
		statusMessage = "Invalid application id"
		response := responses.ApplicationResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	var req requests.UpdateThemeRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		statusMessage = "Invalid request body"
		response := responses.ApplicationResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	if strings.TrimSpace(req.ThemeCode) == "" {
		statusMessage = "theme_code is required"
		response := responses.ApplicationResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	app, err := models.GetApplicationById(id)
	if err != nil {
		statusCode = 404
		statusMessage = "Application not found"
		response := responses.ApplicationResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	theme, err := models.GetThemeByCode(req.ThemeCode)
	if err != nil {
		statusCode = 404
		statusMessage = "Theme not found"
		response := responses.ApplicationResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	appTheme, err := models.GetApplication_themesByApplicationId(id)
	if err != nil {
		statusCode = 404
		statusMessage = "Application theme association not found"
		response := responses.ApplicationResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	appTheme.Theme = theme
	if err := models.UpdateApplication_themesById(appTheme); err != nil {
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
		response := responses.ApplicationResponse{
			StatusCode:    statusCode,
			StatusMessage: statusMessage,
			Result:        &result,
		}
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	updatedAppTheme, err := models.GetApplication_themesByApplicationId(id)
	if err != nil {
		statusCode = 500
		statusMessage = "Theme updated but failed to fetch updated application theme"
	} else {
		statusCode = 200
		statusMessage = "Application theme updated successfully"

		themeConfigs := []responses.ThemePersonalConfigData{}
		for x := range updatedAppTheme.Theme.ThemeConfigs {
			themeConfigs = append(themeConfigs, responses.ThemePersonalConfigData{
				ThemeConfigCode: updatedAppTheme.Theme.ThemeConfigs[x].ThemeConfigCode,
				ThemeProperties: updatedAppTheme.Theme.ThemeConfigs[x].ThemeProperties,
				DateCreated:     updatedAppTheme.Theme.ThemeConfigs[x].DateCreated,
				DateModified:    updatedAppTheme.Theme.ThemeConfigs[x].DateModified,
				CreatedBy:       updatedAppTheme.Theme.ThemeConfigs[x].CreatedBy,
				ModifiedBy:      updatedAppTheme.Theme.ThemeConfigs[x].ModifiedBy,
				Active:          updatedAppTheme.Theme.ThemeConfigs[x].Active,
			})
		}
		result = responses.ApplicationResponseData{
			ApplicationId:    strconv.FormatInt(app.ApplicationId, 10),
			ApplicationCode:  app.ApplicationCode,
			ApplicationName:  app.ApplicationName,
			ApplicationLogo:  app.ApplicationLogo,
			ThemeColors:      app.ThemeColors,
			DefaultFontsize:  app.DefaultFontsize,
			ApplicationImage: app.ApplicationImage,
			DateCreated:      app.DateCreated,
			DateModified:     app.DateModified,
			Active:           app.Active,
			Theme: &responses.ThemePersonalResponseData{
				ThemeId:      strconv.FormatInt(updatedAppTheme.Theme.ThemeId, 10),
				ThemeCode:    updatedAppTheme.Theme.ThemeCode,
				ThemeName:    updatedAppTheme.Theme.ThemeName,
				DateCreated:  updatedAppTheme.Theme.DateCreated,
				DateModified: updatedAppTheme.Theme.DateModified,
				CreatedBy:    updatedAppTheme.Theme.CreatedBy,
				ModifiedBy:   updatedAppTheme.Theme.ModifiedBy,
				ThemeConfig:  themeConfigs,
			},
		}
	}

	response := responses.ApplicationResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        &result,
	}
	logs.Info("UpdateTheme response: ", response)
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
			ApplicationId: strconv.FormatInt(id, 10),
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

// AddApplicationShop ...
// @Title Add Application Shop
// @Description add a shop to the application
// @Param	body		body 	requests.ApplicationShopRequest	true		"body for Shop content"
// @Success 200 {object} responses.ApplicationResponse
// @Failure 403 body is empty
// @router /add-shop [post]
func (c *ApplicationController) AddApplicationShop() {
	var v requests.ApplicationShopRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	statusCode := 400
	statusMessage := "Bad Request"

	appIdInt64, _ := strconv.ParseInt(v.ApplicationId, 10, 64)

	if application, err := models.GetApplicationById(appIdInt64); err == nil {

		if appShop, err := models.GetApplicationShopsByAppAndShop(application.ApplicationCode, v.ShopId); err == nil {
			if appShop != nil {
				logs.Info("Application shop already exists for application code: ", application.ApplicationCode, " and shop ID: ", v.ShopId)
				statusCode = 400
				statusMessage = "Application shop already exists"
				var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
				c.Data["json"] = resp
				c.ServeJSON()
				return
			}
		}
		shopResp := functions.GetShop(&c.Controller, v.ShopId)

		if shopResp.StatusCode == 200 {
			asm := models.ApplicationShops{
				Application: application,
				Shop:        shopResp.Result.ShopId,
				Active:      1,
			}
			if _, err := models.AddApplicationShops(&asm); err == nil {
				message := "Shop Added Successfully"
				statusCode = 200
				statusMessage = message
				appData := responses.ApplicationResponseData{
					ApplicationId:    strconv.FormatInt(application.ApplicationId, 10),
					ApplicationCode:  application.ApplicationCode,
					ApplicationName:  application.ApplicationName,
					ApplicationLogo:  application.ApplicationLogo,
					ThemeColors:      application.ThemeColors,
					DefaultFontsize:  application.DefaultFontsize,
					ApplicationImage: application.ApplicationImage,
					DateCreated:      application.DateCreated,
					DateModified:     application.DateModified,
					Active:           application.Active,
					// Theme: application.Theme,
				}
				var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &appData}
				c.Data["json"] = resp
			} else {
				statusCode = 500
				statusMessage = err.Error()
				var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
				c.Data["json"] = resp
			}
		} else {
			statusCode = shopResp.StatusCode
			statusMessage = shopResp.StatusDesc
			var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
			c.Data["json"] = resp
		}
	} else {
		statusCode = 404
		statusMessage = "Application not found"
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// RemoveApplicationShop ...
// @Title Remove Application Shop
// @Description remove a shop from the application
// @Param	body		body 	requests.ApplicationShopRequest	true		"body for Shop content"
// @Success 200 {object} responses.ApplicationResponse
// @Failure 403 body is empty
// @router /remove-shop [delete]
func (c *ApplicationController) RemoveApplicationShop() {
	logs.Info("Request received to remove shop")
	var v requests.ApplicationShopRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	statusCode := 400
	statusMessage := "Bad Request"

	appIdInt64, _ := strconv.ParseInt(v.ApplicationId, 10, 64)
	if application, err := models.GetApplicationById(appIdInt64); err == nil {
		shopResp := functions.GetShop(&c.Controller, v.ShopId)

		if shopResp.StatusCode == 200 {
			if as, err := models.GetApplicationShopsByAppAndShop(application.ApplicationCode, shopResp.Result.ShopId); err == nil {
				if err := models.DeleteApplicationShops(as.Id); err == nil {
					message := "Shop Removed Successfully"
					statusCode = 200
					statusMessage = message
					appData := responses.ApplicationResponseData{
						ApplicationId:    strconv.FormatInt(application.ApplicationId, 10),
						ApplicationCode:  application.ApplicationCode,
						ApplicationName:  application.ApplicationName,
						ApplicationLogo:  application.ApplicationLogo,
						ThemeColors:      application.ThemeColors,
						DefaultFontsize:  application.DefaultFontsize,
						ApplicationImage: application.ApplicationImage,
						DateCreated:      application.DateCreated,
						DateModified:     application.DateModified,
						Active:           application.Active,
						// Theme: application.Theme,
					}
					var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &appData}
					c.Data["json"] = resp
				} else {
					statusCode = 500
					statusMessage = err.Error()
					var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
					c.Data["json"] = resp
				}
			} else {
				statusCode = 500
				statusMessage = err.Error()
				var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
				c.Data["json"] = resp
			}
		} else {
			statusCode = shopResp.StatusCode
			statusMessage = shopResp.StatusDesc
			var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
			c.Data["json"] = resp
		}
	} else {
		statusCode = 404
		statusMessage = "Application not found"
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// GetApplicationShops ...
// @Title Get Application Shops
// @Description get the shops associated with an application
// @Success 200 {object} responses.ApplicationResponseDTO
// @Failure 400,500 {object} responses.ErrorResponse
// @router /get-application-shops [get]
func (c *ApplicationController) GetApplicationShops() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	logs.Info("GetAll request received")

	statusCode := 400
	statusMessage := "Bad Request"
	var result []responses.ApplicationShopResponseData

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
	query_ := c.GetString("query")
	logs.Info("Query received is ", query_)
	if query_ != "" {
		query_ = strings.TrimSpace(query_) + ",Active:1"
	} else {
		query_ = "Active:1"
	}
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

	l, err := models.GetAllApplicationShops(query, fields, sortby, order, offset, limit)
	if err != nil {
		statusCode = 500
		statusMessage = "Internal server error: " + err.Error()
	} else {
		for _, shop := range l {
			m := shop.(models.ApplicationShops)

			result = append(result, responses.ApplicationShopResponseData{
				ApplicationShopId: strconv.FormatInt(m.Id, 10),
				ShopId:            m.Shop,
				ApplicationId:     strconv.FormatInt(m.Application.ApplicationId, 10),
				ApplicationCode:   m.Application.ApplicationCode,
				ApplicationName:   m.Application.ApplicationName,
				ApplicationLogo:   m.Application.ApplicationLogo,
				ShowBanner:        m.ShowBanner,
				BannerImages:      m.BannerImages,
				BorderRadius:      strconv.FormatFloat(m.BorderRadius, 'f', -1, 64),
				ThemeColors:       m.Application.ThemeColors,
				DefaultFontsize:   m.Application.DefaultFontsize,
				ApplicationImage:  m.Application.ApplicationImage,
				BaseUrl:           m.BaseUrl,
				DateCreated:       m.DateCreated,
				DateModified:      m.DateModified,
				CreatedBy:         m.CreatedBy,
				ModifiedBy:        m.ModifiedBy,
				Active:            m.Active,
			})
		}
		statusCode = 200
		statusMessage = "Shops retrieved successfully"
	}

	resp := responses.ApplicationShopsResponse{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Result:        result,
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

// UploadImage ...
// @Title Upload Image
// @Description upload an image for the application
// @Param	Image		formData 	file	true		"The image file to upload"
// @Param	ItemID		query 	string	true		"The ID of the item to associate the image with"
// @Success 200 {object} models.ItemImageResponseDTO
// @Failure 400,500 {object} models.ErrorResponse
// @router /upload-image [post]
func (c *ApplicationController) UploadImage() {
	// var v models.Item_images

	logs.Info("Data received is ", c.Ctx.Input.Query("ItemID"))

	file, header, err := c.GetFile("Image")
	system := strings.ToLower(c.Ctx.Input.Query("System"))

	logs.Info("Data received is ", file)
	logs.Info("System received is ", system)

	if err != nil {
		// c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Failed to get image file."}
		logs.Error("Failed to get the file ", err)
		c.ServeJSON()
		return
	}
	defer file.Close()

	// Save the uploaded file
	fileName := header.Filename
	logs.Info("File Name Extracted is ", fileName)
	filePath := "/uploads/application/" + system + "/" + fileName // Define your file path
	viewHost, _ := beego.AppConfig.String("imagesBaseUrl")
	viewFilePath := viewHost + filePath
	logs.Info("File Path Extracted is ", filePath)
	host, _ := beego.AppConfig.String("imagesUploadBaseUrl")
	filePath = host + filePath
	logs.Info("Full file path is ", filePath)

	err = c.SaveToFile("Image", filePath)

	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		logs.Error("Error saving file", err)
		// c.Data["json"] = map[string]string{"error": "Failed to save the image file."}
		errorMessage := "Error: Failed to save the image file"

		resp := responses.SystemImageResponseDTO{StatusCode: http.StatusInternalServerError, Result: errorMessage, StatusDesc: "Internal Server Error"}

		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	response := responses.SystemImageResponseDTO{StatusCode: 200, Result: viewFilePath, StatusDesc: "Image uploaded successfully"}
	c.Data["json"] = response
	c.ServeJSON()
}

// UpdateApplicationShopBanner ...
// @Title Update Application Shop Banner
// @Description Update border radius and banner images for a shop
// @Param	body		body 	requests.UpdateApplicationShopBannerRequest	true		"body for banner update"
// @Success 200 {object} responses.ApplicationResponse
// @Failure 400 bad request
// @router /shop/update-banner [put]
func (c *ApplicationController) UpdateApplicationShopBanner() {
	var v requests.UpdateApplicationShopBannerRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	statusCode := 400
	statusMessage := "Bad Request"

	if v.UpdatedBy == "" {
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: 400, StatusMessage: "UpdatedBy is required", Result: nil}
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	if v.ApplicationShopId == "" {
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: "ApplicationShopId is required", Result: nil}
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	appidint, _ := strconv.ParseInt(v.ApplicationShopId, 10, 64)
	if appShop, err := models.GetApplicationShopsById(appidint); err == nil {
		appShop.BannerImages = v.BannerImages
		appShop.BorderRadius = v.BorderRadius
		appShop.ModifiedBy = v.UpdatedBy
		appShop.BaseUrl = v.BaseUrl

		if err := models.UpdateApplicationShopsById(appShop); err == nil {
			statusCode = 200
			statusMessage = "Banner updated successfully"
			appData := responses.ApplicationResponseData{
				ApplicationId:    strconv.FormatInt(appShop.Application.ApplicationId, 10),
				ApplicationCode:  appShop.Application.ApplicationCode,
				ApplicationName:  appShop.Application.ApplicationName,
				ApplicationLogo:  appShop.Application.ApplicationLogo,
				ThemeColors:      appShop.Application.ThemeColors,
				DefaultFontsize:  appShop.Application.DefaultFontsize,
				ApplicationImage: appShop.Application.ApplicationImage,
				DateCreated:      appShop.Application.DateCreated,
				DateModified:     appShop.Application.DateModified,
				Active:           appShop.Application.Active,
			}
			var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &appData}
			c.Data["json"] = resp
		} else {
			statusCode = 500
			statusMessage = err.Error()
			var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
			c.Data["json"] = resp
		}
	} else {
		statusCode = 404
		statusMessage = "Application shop not found"
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// ToggleApplicationShopBanner ...
// @Title Toggle Application Shop Banner
// @Description Toggle the show banner on and off for a shop
// @Param	body		body 	requests.ToggleApplicationShopBannerRequest	true		"body for banner toggle"
// @Success 200 {object} responses.ApplicationResponse
// @Failure 400 bad request
// @router /shop/toggle-banner [put]
func (c *ApplicationController) ToggleApplicationShopBanner() {
	var v requests.ToggleApplicationShopBannerRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if v.UpdatedBy == "" {
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: 400, StatusMessage: "UpdatedBy is required", Result: nil}
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	statusCode := 400
	statusMessage := "Bad Request"

	if v.ApplicationShopId == "" {
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: "ApplicationShopId is required", Result: nil}
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	appidint, _ := strconv.ParseInt(v.ApplicationShopId, 10, 64)
	if appShop, err := models.GetApplicationShopsById(appidint); err == nil {
		appShop.ShowBanner = v.ShowBanner
		appShop.ModifiedBy = v.UpdatedBy

		if err := models.UpdateApplicationShopsById(appShop); err == nil {
			statusCode = 200
			statusMessage = "Banner toggled successfully"
			appData := responses.ApplicationResponseData{
				ApplicationId:    strconv.FormatInt(appShop.Application.ApplicationId, 10),
				ApplicationCode:  appShop.Application.ApplicationCode,
				ApplicationName:  appShop.Application.ApplicationName,
				ApplicationLogo:  appShop.Application.ApplicationLogo,
				ThemeColors:      appShop.Application.ThemeColors,
				DefaultFontsize:  appShop.Application.DefaultFontsize,
				ApplicationImage: appShop.Application.ApplicationImage,
				DateCreated:      appShop.Application.DateCreated,
				DateModified:     appShop.Application.DateModified,
				Active:           appShop.Application.Active,
			}
			var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &appData}
			c.Data["json"] = resp
		} else {
			statusCode = 500
			statusMessage = err.Error()
			var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
			c.Data["json"] = resp
		}
	} else {
		statusCode = 404
		statusMessage = "Application shop not found"
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// UpdateApplicationTheme ...
// @Title Update Application Theme
// @Description Update theme colors and default font size for an application
// @Param	body		body 	requests.UpdateApplicationThemeRequest	true		"body for theme update"
// @Success 200 {object} responses.ApplicationResponse
// @Failure 400 bad request
// @router /update-theme [put]
func (c *ApplicationController) UpdateApplicationTheme() {
	var v requests.UpdateApplicationThemeRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	statusCode := 400
	statusMessage := "Bad Request"

	if v.ApplicationId == "" {
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: "ApplicationId is required", Result: nil}
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	appidint, _ := strconv.ParseInt(v.ApplicationId, 10, 64)
	if app, err := models.GetApplicationById(appidint); err == nil {
		app.ThemeColors = v.ThemeColors
		app.DefaultFontsize = v.DefaultFontsize
		app.ModifiedBy = v.UpdatedBy

		if err := models.UpdateApplicationById(app); err == nil {
			statusCode = 200
			statusMessage = "Application theme updated successfully"
			appData := responses.ApplicationResponseData{
				ApplicationId:    strconv.FormatInt(app.ApplicationId, 10),
				ApplicationCode:  app.ApplicationCode,
				ApplicationName:  app.ApplicationName,
				ApplicationLogo:  app.ApplicationLogo,
				ThemeColors:      app.ThemeColors,
				DefaultFontsize:  app.DefaultFontsize,
				ApplicationImage: app.ApplicationImage,
				DateCreated:      app.DateCreated,
				DateModified:     app.DateModified,
				Active:           app.Active,
			}
			var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: &appData}
			c.Data["json"] = resp
		} else {
			statusCode = 500
			statusMessage = err.Error()
			var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
			c.Data["json"] = resp
		}
	} else {
		statusCode = 404
		statusMessage = "Application not found"
		var resp responses.ApplicationResponse = responses.ApplicationResponse{StatusCode: statusCode, StatusMessage: statusMessage, Result: nil}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}
