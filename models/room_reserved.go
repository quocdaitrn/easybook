package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type RoomReserved struct {
	Id            int          `orm:"column(id);auto"`
	ReservationId *Reservation `orm:"column(reservation_id);rel(fk)"`
	RoomId        *Room        `orm:"column(room_id);rel(fk)"`
	Price         float32      `orm:"column(price);null"`
	CheckIn       time.Time    `orm:"column(check_in);type(timestamp);null"`
	CheckOut      time.Time    `orm:"column(check_out);type(timestamp);null"`
	CreatedAt     time.Time    `orm:"column(created_at);type(timestamp)"`
	UpdatedAt     time.Time    `orm:"column(updated_at);type(timestamp)"`
}

func (t *RoomReserved) TableName() string {
	return "room_reserved"
}

func init() {
	orm.RegisterModel(new(RoomReserved))
}

// AddRoomReserved insert a new RoomReserved into database and returns
// last inserted Id on success.
func AddRoomReserved(m *RoomReserved) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRoomReservedById retrieves RoomReserved by Id. Returns error if
// Id doesn't exist
func GetRoomReservedById(id int) (v *RoomReserved, err error) {
	o := orm.NewOrm()
	v = &RoomReserved{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRoomReserved retrieves all RoomReserved matches certain condition. Returns empty list if
// no records exist
func GetAllRoomReserved(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RoomReserved))
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

	var l []RoomReserved
	qs = qs.OrderBy(sortFields...)
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

// UpdateRoomReserved updates RoomReserved by Id and returns error if
// the record to be updated doesn't exist
func UpdateRoomReservedById(m *RoomReserved) (err error) {
	o := orm.NewOrm()
	v := RoomReserved{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRoomReserved deletes RoomReserved by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRoomReserved(id int) (err error) {
	o := orm.NewOrm()
	v := RoomReserved{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&RoomReserved{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
