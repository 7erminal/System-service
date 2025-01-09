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

// CountriesController operations for Countries
type CountriesController struct {
	beego.Controller
}

// URLMapping ...
func (c *CountriesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Countries
// @Param	body		body 	requests.CountriesRequestDTO	true		"body for Countries content"
// @Success 200 {int} models.Countries
// @Failure 403 body is empty
// @router / [post]
func (c *CountriesController) Post() {
	var v requests.CountriesRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if _, err := models.GetCountriesByCode(v.CountryCode); err != nil {
		currencyId, _ := strconv.ParseInt(v.CurrencyId, 10, 64)
		if currency, err := models.GetCurrenciesById(currencyId); err == nil {
			var country models.Countries = models.Countries{Country: v.Country, CountryCode: v.CountryCode, DefaultCurrency: currency, DateCreated: time.Now(), DateModified: time.Now(), CreatedBy: 1, ModifiedBy: 1}
			if _, err := models.AddCountries(&country); err == nil {
				resp := responses.CountryResponseDTO{StatusCode: 200, Country: &country, StatusDesc: "Country added Successfully"}
				c.Ctx.Output.SetStatus(200)
				c.Data["json"] = resp
			} else {
				logs.Info("Error adding country ", err.Error())
				resp := responses.CountryResponseDTO{StatusCode: 608, Country: nil, StatusDesc: "Error adding country"}
				c.Data["json"] = resp
			}
		} else {
			logs.Info("Error adding country ", err.Error())
			resp := responses.CountryResponseDTO{StatusCode: 608, Country: nil, StatusDesc: "Error adding country"}
			c.Data["json"] = resp
		}
	} else {
		logs.Info("Country already exist")
		resp := responses.CountryResponseDTO{StatusCode: 502, Country: nil, StatusDesc: "Country exists"}
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Countries by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Countries
// @Failure 403 :id is empty
// @router /:id [get]
func (c *CountriesController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetCountriesById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Countries
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Countries
// @Failure 403
// @router / [get]
func (c *CountriesController) GetAll() {
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

	l, err := models.GetAllCountries(query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Info("Error fetching countries ", err.Error())
		resp := responses.CountriesResponseDTO{StatusCode: 608, Countries: nil, StatusDesc: "Error fetching countries"}
		c.Data["json"] = resp
	} else {
		resp := responses.CountriesResponseDTO{StatusCode: 200, Countries: &l, StatusDesc: "Countries fetched successfully"}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Countries
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Countries	true		"body for Countries content"
// @Success 200 {object} models.Countries
// @Failure 403 :id is not int
// @router /:id [put]
func (c *CountriesController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Countries{CountryId: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateCountriesById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Countries
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *CountriesController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteCountries(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
