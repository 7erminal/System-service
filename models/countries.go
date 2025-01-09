package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Countries struct {
	CountryId       int64       `orm:"auto"`
	Country         string      `orm:"size(255)"`
	Description     string      `orm:"size(500)"`
	CountryCode     string      `orm:"size(20)"`
	DefaultCurrency *Currencies `orm:"rel(fk);column(default_currency)"`
	DateCreated     time.Time   `orm:"type(datetime)"`
	DateModified    time.Time   `orm:"type(datetime)"`
	CreatedBy       int
	ModifiedBy      int
}

func init() {
	orm.RegisterModel(new(Countries))
}

func (u *Countries) TableName() string {
	return "countries"
}

// AddCountries insert a new Countries into database and returns
// last inserted Id on success.
func AddCountries(m *Countries) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCountriesById retrieves Countries by Id. Returns error if
// Id doesn't exist
func GetCountriesById(id int64) (v *Countries, err error) {
	o := orm.NewOrm()
	v = &Countries{CountryId: id}
	if err = o.QueryTable(new(Countries)).Filter("CountryId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetCountriesByCode retrieves Countries by Code. Returns error if
// Code doesn't exist
func GetCountriesByCode(countryCode string) (v *Countries, err error) {
	o := orm.NewOrm()
	v = &Countries{CountryCode: countryCode}
	if err = o.QueryTable(new(Countries)).Filter("CountryCode", countryCode).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCountries retrieves all Countries matches certain condition. Returns empty list if
// no records exist
func GetAllCountries(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Countries))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Countries
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateCountries updates Countries by Id and returns error if
// the record to be updated doesn't exist
func UpdateCountriesById(m *Countries) (err error) {
	o := orm.NewOrm()
	v := Countries{CountryId: m.CountryId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCountries deletes Countries by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCountries(id int64) (err error) {
	o := orm.NewOrm()
	v := Countries{CountryId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Countries{CountryId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
