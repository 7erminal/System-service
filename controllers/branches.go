package controllers

import (
	"encoding/json"
	"errors"
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

// BranchesController operations for Branches
type BranchesController struct {
	beego.Controller
}

// URLMapping ...
func (c *BranchesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Branches
// @Param	body		body 	requests.BranchRequestDTO	true		"body for Branches content"
// @Success 201 {int} models.Branches
// @Failure 403 body is empty
// @router / [post]
func (c *BranchesController) Post() {
	var v requests.BranchRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	logs.Info("User ID searched is ", v.AddedBy)

	userid, _ := strconv.ParseInt(v.AddedBy, 10, 64)

	userCheck := functions.GetUserDetails(&c.Controller, userid)

	if userCheck.StatusCode == 200 {
		if country, err := models.GetCountriesByCode(v.CountryCode); err == nil {
			var branch models.Branches = models.Branches{Branch: v.Branch, Country: country, Location: v.Location, PhoneNumber: v.PhoneNumber, Active: 1, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: int(userid), ModifiedBy: int(userid)}
			if _, err := models.AddBranches(&branch); err == nil {
				resp := responses.BranchResponseDTO{StatusCode: 200, Branch: &branch, StatusDesc: "Branch added successfully"}
				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = resp
			} else {
				logs.Info("Error adding branch ", err.Error())
				resp := responses.BranchResponseDTO{StatusCode: 402, Branch: nil, StatusDesc: "Branch additon failed"}
				c.Data["json"] = resp
			}
		} else {
			logs.Info("Error adding branch ", err.Error())
			resp := responses.BranchResponseDTO{StatusCode: 608, Branch: nil, StatusDesc: "Error getting specified country"}
			c.Data["json"] = resp
		}
	} else {
		logs.Info("Error::: User not found ")
		resp := responses.BranchResponseDTO{StatusCode: 500, Branch: nil, StatusDesc: "Error adding Branch"}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Branches by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Branches
// @Failure 403 :id is empty
// @router /:id [get]
func (c *BranchesController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	logs.Info("Getting branch details:: ", id)
	v, err := models.GetBranchesById(id)
	if err != nil {
		logs.Error("Error getting branch ", err.Error())
		resp := responses.BranchResponseDTO{StatusCode: 402, Branch: nil, StatusDesc: "Error getting branch"}
		c.Data["json"] = resp
	} else {
		resp := responses.BranchResponseDTO{StatusCode: 200, Branch: v, StatusDesc: "Branch details fetched successfully"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Branches
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Branches
// @Failure 403
// @router / [get]
func (c *BranchesController) GetAll() {
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

	l, err := models.GetAllBranches(query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Info("Error fetching branches ", err.Error())
		resp := responses.BranchesResponseDTO{StatusCode: 608, Branches: nil, StatusDesc: "Error fetching branches"}
		c.Data["json"] = resp
	} else {
		resp := responses.BranchesResponseDTO{StatusCode: 200, Branches: &l, StatusDesc: "Branches fetched successfully"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Branches
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Branches	true		"body for Branches content"
// @Success 200 {object} responses.BranchResponseDTO
// @Failure 403 :id is not int
// @router /:id [put]
func (c *BranchesController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Branches{BranchId: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateBranchesById(&v); err == nil {
		resp := responses.BranchResponseDTO{StatusCode: 200, Branch: &v, StatusDesc: "Branch updated successfully"}
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = resp
	} else {
		logs.Error("Branch update failed", err.Error())
		resp := responses.BranchResponseDTO{StatusCode: 608, Branch: nil, StatusDesc: "Branch update failed"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// UpdateBranchManager ...
// @Title Put Branch Manager
// @Description update the Branches manager
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Branches	true		"body for Branches content"
// @Success 200 {object} responses.BranchResponseDTO
// @Failure 403 :id is not int
// @router /branch-manager/:id [put]
func (c *BranchesController) UpdateBranchManager() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	var v requests.BranchManagerRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	userId, _ := strconv.ParseInt(v.BranchManager, 0, 64)

	if user, err := models.GetUsersById(userId); err == nil {

		if branch, err := models.GetBranchesById(id); err == nil {
			branch.BranchManager = user
			if err := models.UpdateBranchesById(branch); err == nil {
				resp := responses.BranchResponseDTO{StatusCode: 200, Branch: branch, StatusDesc: "Branch updated successfully"}
				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = resp
			} else {
				logs.Error("Branch update failed", err.Error())
				resp := responses.BranchResponseDTO{StatusCode: 608, Branch: nil, StatusDesc: "Branch update failed"}
				c.Data["json"] = resp
			}
		} else {
			logs.Error("Branch update failed", err.Error())
			resp := responses.BranchResponseDTO{StatusCode: 608, Branch: nil, StatusDesc: "Branch update failed. Branch not found"}
			c.Data["json"] = resp
		}

	} else {
		logs.Error("Branch update failed", err.Error())
		resp := responses.BranchResponseDTO{StatusCode: 608, Branch: nil, StatusDesc: "Branch update failed. User not found"}
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Branches
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *BranchesController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteBranches(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
