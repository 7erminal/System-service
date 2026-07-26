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

	beego "github.com/beego/beego/v2/server/web"
)

// BillersController operations for Billers
type BillersController struct {
	beego.Controller
}

// URLMapping ...
func (c *BillersController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Billers
// @Param	body		body 	requests.BillerRequestDTO	true		"body for Billers content"
// @Success 201 {int} responses.BillerResponseDTO
// @Failure 403 body is empty
// @router / [post]
func (c *BillersController) Post() {
	var req requests.BillerRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	if operator, err := models.GetOperatorByCode(req.OperatorCode); err == nil {
		biller := models.Billers{
			BillerName:   req.BillerName,
			CreatedBy:    req.CreatedBy,
			ModifiedBy:   req.CreatedBy,
			DateCreated:  time.Now(),
			DateModified: time.Now(),
			Active:       1,
			Operator:     operator,
		}

		if _, err := models.AddBillers(&biller); err == nil {
			c.Ctx.Output.SetStatus(201)

			dateCreatedStr := biller.DateCreated.Format("2006-01-02 15:04:05")
			dateModifiedStr := biller.DateModified.Format("2006-01-02 15:04:05")
			billerObj := responses.BillerObject{
				BillerId:     biller.BillerId,
				BillerName:   biller.BillerName,
				CreatedBy:    biller.CreatedBy,
				ModifiedBy:   biller.ModifiedBy,
				DateCreated:  dateCreatedStr,
				DateModified: dateModifiedStr,
				Active:       biller.Active,
				Operator:     biller.Operator.OperatorCode,
			}
			billerResp := responses.BillerResponseDTO{
				StatusCode: 200,
				Biller:     &billerObj,
				StatusDesc: "Biller added successfully",
			}
			c.Data["json"] = billerResp
		} else {
			billerResp := responses.BillerResponseDTO{
				StatusCode: 500,
				Biller:     nil,
				StatusDesc: "Failed to add Biller: " + err.Error(),
			}
			c.Data["json"] = billerResp
		}
	} else {
		billerResp := responses.BillerResponseDTO{
			StatusCode: 500,
			Biller:     nil,
			StatusDesc: "Failed to add Biller: " + err.Error(),
		}
		c.Data["json"] = billerResp
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Billers by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Billers
// @Failure 403 :id is empty
// @router /:id [get]
func (c *BillersController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetBillersById(id)
	if err != nil {
		billerResp := responses.BillerResponseDTO{
			StatusCode: 500,
			Biller:     nil,
			StatusDesc: "Failed to fetch Biller: " + err.Error(),
		}
		c.Data["json"] = billerResp
	} else {
		dateCreatedStr := v.DateCreated.Format("2006-01-02 15:04:05")
		dateModifiedStr := v.DateModified.Format("2006-01-02 15:04:05")
		billerObj := responses.BillerObject{
			BillerId:     v.BillerId,
			BillerName:   v.BillerName,
			CreatedBy:    v.CreatedBy,
			ModifiedBy:   v.ModifiedBy,
			DateCreated:  dateCreatedStr,
			DateModified: dateModifiedStr,
			Active:       v.Active,
			Operator:     v.Operator.OperatorCode,
		}
		billerResp := responses.BillerResponseDTO{
			StatusCode: 200,
			Biller:     &billerObj,
			StatusDesc: "Biller fetched successfully",
		}
		c.Data["json"] = billerResp
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Billers
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Billers
// @Failure 403
// @router / [get]
func (c *BillersController) GetAll() {
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

	l, err := models.GetAllBillers(query, fields, sortby, order, offset, limit)
	if err != nil {
		billersResp := responses.BillersResponseDTO{
			StatusCode: 500,
			Result:     []responses.BillerObject{},
			StatusDesc: err.Error(),
		}
		c.Data["json"] = billersResp
	} else {
		billerObj := []responses.BillerObject{}
		for _, urs := range l {
			m := urs.(models.Billers)
			dateCreatedStr := m.DateCreated.Format("2006-01-02 15:04:05")
			dateModifiedStr := m.DateModified.Format("2006-01-02 15:04:05")
			m.DateCreated, _ = time.Parse("2006-01-02 15:04:05", dateCreatedStr)
			m.DateModified, _ = time.Parse("2006-01-02 15:04:05", dateModifiedStr)
			billerObj = append(billerObj, responses.BillerObject{
				BillerId:     m.BillerId,
				BillerName:   m.BillerName,
				CreatedBy:    m.CreatedBy,
				ModifiedBy:   m.ModifiedBy,
				DateCreated:  dateCreatedStr,
				DateModified: dateModifiedStr,
				Active:       m.Active,
				Operator:     m.Operator.OperatorCode,
			})
		}

		billersResp := responses.BillersResponseDTO{
			StatusCode: 200,
			Result:     billerObj,
			StatusDesc: "Billers fetched successfully",
		}
		c.Data["json"] = billersResp
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Billers
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Billers	true		"body for Billers content"
// @Success 200 {object} models.Billers
// @Failure 403 :id is not int
// @router /:id [put]
func (c *BillersController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Billers{BillerId: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateBillersById(&v); err == nil {
		billerResp := responses.BillerResponseDTO{
			StatusCode: 200,
			Biller:     nil,
			StatusDesc: "Biller updated successfully",
		}
		c.Data["json"] = billerResp
	} else {
		billerResp := responses.BillerResponseDTO{
			StatusCode: 500,
			Biller:     nil,
			StatusDesc: "Failed to update Biller: " + err.Error(),
		}
		c.Data["json"] = billerResp
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Billers
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *BillersController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteBillers(id); err == nil {
		billerResp := responses.BillerResponseDTO{
			StatusCode: 200,
			Biller:     nil,
			StatusDesc: "Biller deleted successfully",
		}
		c.Data["json"] = billerResp
	} else {
		billerResp := responses.BillerResponseDTO{
			StatusCode: 500,
			Biller:     nil,
			StatusDesc: "Failed to delete Biller: " + err.Error(),
		}
		c.Data["json"] = billerResp
	}
	c.ServeJSON()
}
