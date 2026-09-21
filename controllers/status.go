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
// @Param	body		body 	requests.Status	true		"body for Status content"
// @Success 201 {int} requests.Status
// @Failure 403 body is empty
// @router / [post]
func (c *StatusController) Post() {
	var req requests.Status
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	statusCode := 608
	message := "CreateStatus error"

	v := models.Status{
		Status:     req.Status,
		StatusCode: req.StatusCode,
		Active:     1,
	}
	if _, err := models.AddStatus(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		statusCode = 200
		message = "CreateStatus success"
		resp := responses.StatusResponse{
			StatusCode: statusCode,
			StatusDesc: message,
			Result: &responses.StatusResponseDTO{
				Status:     v.Status,
				StatusCode: v.StatusCode,
				StatusId:   v.StatusId,
			},
		}
		c.Data["json"] = resp
	} else {
		c.Data["json"] = responses.StatusResponse{
			StatusCode: statusCode,
			StatusDesc: message,
			Result:     nil,
		}
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
	statusCode := 608
	message := "GetStatusById error"
	resp := responses.StatusResponseDTO{}
	if err != nil {
		logs.Error(err)
		statusCode = 608
		message = "GetStatusById error"
		c.Data["json"] = responses.StatusResponse{
			StatusCode: statusCode,
			StatusDesc: message,
			Result:     nil,
		}
	} else {
		statusCode = 200
		message = "GetStatusById success"
		resp = responses.StatusResponseDTO{
			Status:     v.Status,
			StatusCode: v.StatusCode,
			StatusId:   v.StatusId,
		}
		response := responses.StatusResponse{
			StatusCode: statusCode,
			StatusDesc: message,
			Result:     &resp,
		}
		c.Data["json"] = response
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
		resp := responses.StatusListResponse{StatusCode: statusCode, StatusDesc: message, Result: nil}
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
		resp := responses.StatusListResponse{StatusCode: statusCode, StatusDesc: message, Result: &statusDTOs}
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
	statusCode := 608
	message := "UpdateStatus error"
	resp := responses.StatusResponseDTO{}
	req := requests.Status{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	v := models.Status{
		StatusId:   id,
		StatusCode: req.StatusCode,
		Status:     req.Status,
	}
	if err := models.UpdateStatusById(&v); err == nil {
		statusCode = 200
		message = "UpdateStatus success"
		resp = responses.StatusResponseDTO{
			StatusId:   v.StatusId,
			StatusCode: v.StatusCode,
			Status:     v.Status,
		}
		c.Data["json"] = responses.StatusResponse{
			StatusCode: statusCode,
			StatusDesc: message,
			Result:     &resp,
		}
	} else {
		logs.Error("UpdateStatus error: ", err)
		c.Data["json"] = responses.StatusResponse{
			StatusCode: statusCode,
			StatusDesc: message,
			Result:     nil,
		}
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
		statusCode := 200
		message := "DeleteStatus success"
		c.Data["json"] = responses.StatusResponse{
			StatusCode: statusCode,
			StatusDesc: message,
			Result:     nil,
		}
	} else {
		statusCode := 608
		message := "DeleteStatus error"
		logs.Error("DeleteStatus error: ", err)
		c.Data["json"] = responses.StatusResponse{
			StatusCode: statusCode,
			StatusDesc: message,
			Result:     nil,
		}
	}
	c.ServeJSON()
}
