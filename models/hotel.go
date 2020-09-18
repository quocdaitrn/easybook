package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Hotel struct {
	Id          int       `orm:"column(id);auto"`
	Name        string    `orm:"column(name);size(100)"`
	Description string    `orm:"column(description);null"`
	IsActive    int8      `orm:"column(is_active)"`
	Address     string    `orm:"column(address);size(255)"`
	CityId      *City     `orm:"column(city_id);rel(fk)"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp)"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp)"`
}

func (t *Hotel) TableName() string {
	return "hotel"
}

func init() {
	orm.RegisterModel(new(Hotel))
}

// AddHotel insert a new Hotel into database and returns
// last inserted Id on success.
func AddHotel(m *Hotel) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetHotelById retrieves Hotel by Id. Returns error if
// Id doesn't exist
func GetHotelById(id int) (v *Hotel, err error) {
	o := orm.NewOrm()
	v = &Hotel{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllHotel retrieves all Hotel matches certain condition. Returns empty list if
// no records exist
func GetAllHotel(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Hotel))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
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

	var l []Hotel
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.RelatedSel().Limit(limit, offset).All(&l, fields...); err == nil {
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

// UpdateHotel updates Hotel by Id and returns error if
// the record to be updated doesn't exist
func UpdateHotelById(m *Hotel) (err error) {
	o := orm.NewOrm()
	v := Hotel{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteHotel deletes Hotel by Id and returns error if
// the record to be deleted doesn't exist
func DeleteHotel(id int) (err error) {
	o := orm.NewOrm()
	v := Hotel{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Hotel{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
