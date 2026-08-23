package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Application_themes struct {
	Id            int64       `orm:"auto;column(application_theme_id)"`
	ApplicationId Application `orm:"rel(fk)"`
	ThemeId       Theme       `orm:"rel(fk)"`
	DateCreated   time.Time   `orm:"type(datetime)"`
	DateModified  time.Time   `orm:"type(datetime)"`
	CreatedBy     int
	ModifiedBy    int
	Active        int
}

func init() {
	orm.RegisterModel(new(Application_themes))
}

// AddApplication_themes insert a new Application_themes into database and returns
// last inserted Id on success.
func AddApplication_themes(m *Application_themes) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetApplication_themesById retrieves Application_themes by Id. Returns error if
// Id doesn't exist
func GetApplication_themesById(id int64) (v *Application_themes, err error) {
	o := orm.NewOrm()
	v = &Application_themes{Id: id}
	if err = o.QueryTable(new(Application_themes)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetApplication_themesByApplicationId(id int64) (v *Application_themes, err error) {
	o := orm.NewOrm()
	application := Application{ApplicationId: id}
	v = &Application_themes{ApplicationId: application}
	if err = o.QueryTable(new(Application_themes)).Filter("ApplicationId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApplication_themes retrieves all Application_themes matches certain condition. Returns empty list if
// no records exist
func GetAllApplication_themes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Application_themes))
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

	var l []Application_themes
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

// UpdateApplication_themes updates Application_themes by Id and returns error if
// the record to be updated doesn't exist
func UpdateApplication_themesById(m *Application_themes) (err error) {
	o := orm.NewOrm()
	v := Application_themes{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApplication_themes deletes Application_themes by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApplication_themes(id int64) (err error) {
	o := orm.NewOrm()
	v := Application_themes{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Application_themes{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
