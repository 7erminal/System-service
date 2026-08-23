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

// OperatorController operations for Operator
type OperatorController struct {
	beego.Controller
}

// URLMapping ...
func (c *OperatorController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Operator
// @Param	body		body 	models.Operator	true		"body for Operator content"
// @Success 201 {int} models.Operator
// @Failure 403 body is empty
// @router / [post]
func (c *OperatorController) Post() {
	var v models.Operator
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 400,
			Operator:   nil,
			StatusDesc: "Invalid request body: " + err.Error(),
		}
		c.Data["json"] = operatorResp
		c.ServeJSON()
		return
	}
	if _, err := models.AddOperator(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		dateCreatedStr := v.DateCreated.Format("2006-01-02 15:04:05")
		dateModifiedStr := v.DateModified.Format("2006-01-02 15:04:05")
		operatorObj := responses.OperatorObject{
			OperatorId:   v.OperatorId,
			OperatorName: v.OperatorName,
			Description:  v.Description,
			DateCreated:  dateCreatedStr,
			DateModified: dateModifiedStr,
			CreatedBy:    v.CreatedBy,
			ModifiedBy:   v.ModifiedBy,
			Active:       v.Active,
		}
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 200,
			Operator:   &operatorObj,
			StatusDesc: "Operator added successfully",
		}
		c.Data["json"] = operatorResp
	} else {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 500,
			Operator:   nil,
			StatusDesc: "Failed to add Operator: " + err.Error(),
		}
		c.Data["json"] = operatorResp
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Operator by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Operator
// @Failure 403 :id is empty
// @router /:id [get]
func (c *OperatorController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 400,
			Operator:   nil,
			StatusDesc: "Invalid operator id",
		}
		c.Data["json"] = operatorResp
		c.ServeJSON()
		return
	}
	v, err := models.GetOperatorById(id)
	if err != nil {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 500,
			Operator:   nil,
			StatusDesc: "Failed to fetch Operator: " + err.Error(),
		}
		c.Data["json"] = operatorResp
	} else {
		dateCreatedStr := v.DateCreated.Format("2006-01-02 15:04:05")
		dateModifiedStr := v.DateModified.Format("2006-01-02 15:04:05")
		operatorObj := responses.OperatorObject{
			OperatorId:   v.OperatorId,
			OperatorName: v.OperatorName,
			Description:  v.Description,
			DateCreated:  dateCreatedStr,
			DateModified: dateModifiedStr,
			CreatedBy:    v.CreatedBy,
			ModifiedBy:   v.ModifiedBy,
			Active:       v.Active,
		}
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 200,
			Operator:   &operatorObj,
			StatusDesc: "Operator fetched successfully",
		}
		c.Data["json"] = operatorResp
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Operator
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Operator
// @Failure 403
// @router / [get]
func (c *OperatorController) GetAll() {
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
				operatorsResp := responses.OperatorsResponseDTO{
					StatusCode: 400,
					Operator:   []responses.OperatorObject{},
					StatusDesc: errors.New("Error: invalid query key/value pair").Error(),
				}
				c.Data["json"] = operatorsResp
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllOperator(query, fields, sortby, order, offset, limit)
	if err != nil {
		operatorsResp := responses.OperatorsResponseDTO{
			StatusCode: 500,
			Operator:   []responses.OperatorObject{},
			StatusDesc: "Failed to fetch Operators: " + err.Error(),
		}
		c.Data["json"] = operatorsResp
	} else {
		operatorObj := []responses.OperatorObject{}
		for _, urs := range l {
			m := urs.(models.Operator)
			dateCreatedStr := m.DateCreated.Format("2006-01-02 15:04:05")
			dateModifiedStr := m.DateModified.Format("2006-01-02 15:04:05")
			operatorObj = append(operatorObj, responses.OperatorObject{
				OperatorId:   m.OperatorId,
				OperatorName: m.OperatorName,
				Description:  m.Description,
				DateCreated:  dateCreatedStr,
				DateModified: dateModifiedStr,
				CreatedBy:    m.CreatedBy,
				ModifiedBy:   m.ModifiedBy,
				Active:       m.Active,
			})
		}

		operatorsResp := responses.OperatorsResponseDTO{
			StatusCode: 200,
			Operator:   operatorObj,
			StatusDesc: "Operators fetched successfully",
		}
		c.Data["json"] = operatorsResp
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Operator
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Operator	true		"body for Operator content"
// @Success 200 {object} models.Operator
// @Failure 403 :id is not int
// @router /:id [put]
func (c *OperatorController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 400,
			Operator:   nil,
			StatusDesc: "Invalid operator id",
		}
		c.Data["json"] = operatorResp
		c.ServeJSON()
		return
	}
	v := models.Operator{OperatorId: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 400,
			Operator:   nil,
			StatusDesc: "Invalid request body: " + err.Error(),
		}
		c.Data["json"] = operatorResp
		c.ServeJSON()
		return
	}
	if err := models.UpdateOperatorById(&v); err == nil {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 200,
			Operator:   nil,
			StatusDesc: "Operator updated successfully",
		}
		c.Data["json"] = operatorResp
	} else {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 500,
			Operator:   nil,
			StatusDesc: "Failed to update Operator: " + err.Error(),
		}
		c.Data["json"] = operatorResp
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Operator
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *OperatorController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 400,
			Operator:   nil,
			StatusDesc: "Invalid operator id",
		}
		c.Data["json"] = operatorResp
		c.ServeJSON()
		return
	}
	if err := models.DeleteOperator(id); err == nil {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 200,
			Operator:   nil,
			StatusDesc: "Operator deleted successfully",
		}
		c.Data["json"] = operatorResp
	} else {
		operatorResp := responses.OperatorResponseDTO{
			StatusCode: 500,
			Operator:   nil,
			StatusDesc: "Failed to delete Operator: " + err.Error(),
		}
		c.Data["json"] = operatorResp
	}
	c.ServeJSON()
}
