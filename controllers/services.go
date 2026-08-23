package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"system_service/models"
	"system_service/structs/responses"

	beego "github.com/beego/beego/v2/server/web"
)

// ServicesController operations for Services
type ServicesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ServicesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Services
// @Param	body		body 	models.Services	true		"body for Services content"
// @Success 201 {int} models.Services
// @Failure 403 body is empty
// @router / [post]
func (c *ServicesController) Post() {
	var v models.Services
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 400,
			Result:     responses.ServiceObject{},
			StatusDesc: "Invalid request body: " + err.Error(),
		}
		c.Data["json"] = serviceResp
		c.ServeJSON()
		return
	}
	if _, err := models.AddServices(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		dateCreatedStr := v.DateCreated.Format("2006-01-02 15:04:05")
		dateModifiedStr := v.DateModified.Format("2006-01-02 15:04:05")
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 200,
			Result: responses.ServiceObject{
				ServiceId:          v.ServiceId,
				ServiceName:        v.ServiceName,
				ServiceCode:        v.ServiceCode,
				ServiceDescription: v.ServiceDescription,
				DateCreated:        dateCreatedStr,
				DateModified:       dateModifiedStr,
				CreatedBy:          v.CreatedBy,
				ModifiedBy:         v.ModifiedBy,
				Active:             v.Active,
			},
			StatusDesc: "Service added successfully",
		}
		c.Data["json"] = serviceResp
	} else {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 500,
			Result:     responses.ServiceObject{},
			StatusDesc: "Failed to add Service: " + err.Error(),
		}
		c.Data["json"] = serviceResp
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Services by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Services
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ServicesController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 400,
			Result:     responses.ServiceObject{},
			StatusDesc: "Invalid service id",
		}
		c.Data["json"] = serviceResp
		c.ServeJSON()
		return
	}
	v, err := models.GetServicesById(id)
	if err != nil {
		v, err := models.GetServicesByCode(idStr)
		if err != nil {
			serviceResp := responses.ServiceResponseDTO{
				StatusCode: 500,
				Result:     responses.ServiceObject{},
				StatusDesc: "Failed to fetch Service: " + err.Error(),
			}
			c.Data["json"] = serviceResp
		} else {
			dateCreatedStr := v.DateCreated.Format("2006-01-02 15:04:05")
			dateModifiedStr := v.DateModified.Format("2006-01-02 15:04:05")
			serviceResp := responses.ServiceResponseDTO{
				StatusCode: 200,
				Result: responses.ServiceObject{
					ServiceId:          v.ServiceId,
					ServiceName:        v.ServiceName,
					ServiceCode:        v.ServiceCode,
					ServiceDescription: v.ServiceDescription,
					DateCreated:        dateCreatedStr,
					DateModified:       dateModifiedStr,
					CreatedBy:          v.CreatedBy,
					ModifiedBy:         v.ModifiedBy,
					Active:             v.Active,
				},
				StatusDesc: "Service fetched successfully",
			}
			c.Data["json"] = serviceResp
		}
	} else {
		dateCreatedStr := v.DateCreated.Format("2006-01-02 15:04:05")
		dateModifiedStr := v.DateModified.Format("2006-01-02 15:04:05")
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 200,
			Result: responses.ServiceObject{
				ServiceId:          v.ServiceId,
				ServiceName:        v.ServiceName,
				ServiceCode:        v.ServiceCode,
				ServiceDescription: v.ServiceDescription,
				DateCreated:        dateCreatedStr,
				DateModified:       dateModifiedStr,
				CreatedBy:          v.CreatedBy,
				ModifiedBy:         v.ModifiedBy,
				Active:             v.Active,
			},
			StatusDesc: "Service fetched successfully",
		}
		c.Data["json"] = serviceResp
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Services
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Services
// @Failure 403
// @router / [get]
func (c *ServicesController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

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
				servicesResp := responses.ServicesResponseDTO{
					StatusCode: 400,
					Result:     []responses.ServiceObject{},
					StatusDesc: errors.New("Error: invalid query key/value pair").Error(),
				}
				c.Data["json"] = servicesResp
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllServices(query, fields, sortby, order, offset, limit)
	if err != nil {
		servicesResp := responses.ServicesResponseDTO{
			StatusCode: 500,
			Result:     []responses.ServiceObject{},
			StatusDesc: "Failed to fetch Services: " + err.Error(),
		}
		c.Data["json"] = servicesResp
	} else {
		serviceObj := []responses.ServiceObject{}
		for _, urs := range l {
			m := urs.(models.Services)
			dateCreatedStr := m.DateCreated.Format("2006-01-02 15:04:05")
			dateModifiedStr := m.DateModified.Format("2006-01-02 15:04:05")
			serviceObj = append(serviceObj, responses.ServiceObject{
				ServiceId:          m.ServiceId,
				ServiceName:        m.ServiceName,
				ServiceCode:        m.ServiceCode,
				ServiceDescription: m.ServiceDescription,
				DateCreated:        dateCreatedStr,
				DateModified:       dateModifiedStr,
				CreatedBy:          m.CreatedBy,
				ModifiedBy:         m.ModifiedBy,
				Active:             m.Active,
			})
		}

		servicesResp := responses.ServicesResponseDTO{
			StatusCode: 200,
			Result:     serviceObj,
			StatusDesc: "Services fetched successfully",
		}
		c.Data["json"] = servicesResp
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Services
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Services	true		"body for Services content"
// @Success 200 {object} models.Services
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ServicesController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 400,
			Result:     responses.ServiceObject{},
			StatusDesc: "Invalid service id",
		}
		c.Data["json"] = serviceResp
		c.ServeJSON()
		return
	}
	v := models.Services{ServiceId: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 400,
			Result:     responses.ServiceObject{},
			StatusDesc: "Invalid request body: " + err.Error(),
		}
		c.Data["json"] = serviceResp
		c.ServeJSON()
		return
	}
	if err := models.UpdateServicesById(&v); err == nil {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 200,
			Result:     responses.ServiceObject{},
			StatusDesc: "Service updated successfully",
		}
		c.Data["json"] = serviceResp
	} else {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 500,
			Result:     responses.ServiceObject{},
			StatusDesc: "Failed to update Service: " + err.Error(),
		}
		c.Data["json"] = serviceResp
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Services
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ServicesController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 400,
			Result:     responses.ServiceObject{},
			StatusDesc: "Invalid service id",
		}
		c.Data["json"] = serviceResp
		c.ServeJSON()
		return
	}
	if err := models.DeleteServices(id); err == nil {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 200,
			Result:     responses.ServiceObject{},
			StatusDesc: "Service deleted successfully",
		}
		c.Data["json"] = serviceResp
	} else {
		serviceResp := responses.ServiceResponseDTO{
			StatusCode: 500,
			Result:     responses.ServiceObject{},
			StatusDesc: "Failed to delete Service: " + err.Error(),
		}
		c.Data["json"] = serviceResp
	}
	c.ServeJSON()
}
