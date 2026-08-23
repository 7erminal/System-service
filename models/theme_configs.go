package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Theme_configs struct {
	Id              int64     `orm:"auto"`
	ThemeId         *Theme    `orm:"rel(fk)"`
	ThemeConfigCode string    `orm:"size(255)"`
	ThemeProperties string    `orm:"type(longtext)"`
	DateCreated     time.Time `orm:"type(datetime)"`
	DateModified    time.Time `orm:"type(datetime)"`
	CreatedBy       int
	ModifiedBy      int
	Active          int
}

func init() {
	orm.RegisterModel(new(Theme_configs))
}

// AddTheme_configs insert a new Theme_configs into database and returns
// last inserted Id on success.
func AddTheme_configs(m *Theme_configs) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTheme_configsById retrieves Theme_configs by Id. Returns error if
// Id doesn't exist
func GetTheme_configsById(id int64) (v *Theme_configs, err error) {
	o := orm.NewOrm()
	v = &Theme_configs{Id: id}
	if err = o.QueryTable(new(Theme_configs)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetTheme_configsByCode(code string) (v *Theme_configs, err error) {
	o := orm.NewOrm()
	v = &Theme_configs{ThemeConfigCode: code}
	if err = o.QueryTable(new(Theme_configs)).Filter("ThemeConfigCode", code).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetTheme_configsByThemeId(id int64) (v *Theme_configs, err error) {
	o := orm.NewOrm()

	v = &Theme_configs{ThemeId: &Theme{ThemeId: id}}
	if err = o.QueryTable(new(Theme_configs)).Filter("ThemeId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTheme_configs retrieves all Theme_configs matches certain condition. Returns empty list if
// no records exist
func GetAllTheme_configs(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Theme_configs))
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

	var l []Theme_configs
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
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

// UpdateTheme_configs updates Theme_configs by Id and returns error if
// the record to be updated doesn't exist
func UpdateTheme_configsById(m *Theme_configs) (err error) {
	o := orm.NewOrm()
	v := Theme_configs{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTheme_configs deletes Theme_configs by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTheme_configs(id int64) (err error) {
	o := orm.NewOrm()
	v := Theme_configs{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Theme_configs{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
