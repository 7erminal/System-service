package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"system_service/models"
	"system_service/structs/responses"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// StatusController operations for Status
type StatusController struct {
	beego.Controller
}

// URLMapping ...
func (c *StatusController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Status
// @Param	body		body 	models.Status	true		"body for Status content"
// @Success 201 {int} models.Status
// @Failure 403 body is empty
// @router / [post]
func (c *StatusController) Post() {
	var v models.Status
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if _, err := models.AddStatus(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = v
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Status by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Status
// @Failure 403 :id is empty
// @router /:id [get]
func (c *StatusController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetStatusById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Status
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} responses.StatusListResponse
// @Failure 403
// @router / [get]
func (c *StatusController) GetAll() {
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
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	statusCode := 608
	message := "GetAllStatus error"

	l, err := models.GetAllStatus(query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Error("GetAllStatus error: ", err)
		resp := responses.StatusListResponse{StatusCode: statusCode, StatusDesc: message, Statuses: nil}
		c.Data["json"] = resp
		// c.Data["json"] = err.Error()
	} else {
		statusCode = 200
		message = "GetAllStatus success"
		var statuses []models.Status
		for _, item := range l {
			if status, ok := item.(models.Status); ok {
				statuses = append(statuses, status)
			}
		}

		var statusDTOs []responses.StatusResponseDTO
		for _, status := range statuses {
			statusDTO := responses.StatusResponseDTO{
				StatusId:   status.StatusId,
				StatusCode: status.StatusCode,
				Status:     status.Status,
			}
			statusDTOs = append(statusDTOs, statusDTO)
		}
		resp := responses.StatusListResponse{StatusCode: statusCode, StatusDesc: message, Statuses: &statusDTOs}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Status
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Status	true		"body for Status content"
// @Success 200 {object} models.Status
// @Failure 403 :id is not int
// @router /:id [put]
func (c *StatusController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Status{StatusId: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateStatusById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Status
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *StatusController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteStatus(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
