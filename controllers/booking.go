package controllers

import (
	"easybook/reqres"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"easybook/models"
	"easybook/services/easybook_chaincode"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// BookingController operations for booking process
type BookingController struct {
	beego.Controller
}

// URLMapping ...
func (c *BookingController) URLMapping() {
	c.Mapping("SearchHotels", c.SearchHotels)
	c.Mapping("ReserveRooms", c.ReserveRooms)
}

// SearchHotels ...
// @Title Search Hotels
// @Description search available Hotel
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Hotel
// @Failure 403
// @router /search [get]
func (c *BookingController) SearchHotels() {
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

	l, err := models.GetAllHotel(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	contract, err := easybook_chaincode.GetContract()
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	hotels := []*models.Hotel{}
	for _, h := range l {
		hotel := h.(models.Hotel)
		smJsonHotel, err := contract.EvaluateTransaction("ReadHotel", strconv.Itoa(hotel.Id))
		if err == nil {
			smHotel := &easybook_chaincode.Hotel{}
			err = json.Unmarshal(smJsonHotel, smHotel)
			if err != nil {
				c.Data["json"] = err.Error()
				c.ServeJSON()
				return
			}
			hotel.Rating = smHotel.Rating
		}
		hotels = append(hotels, &hotel)
	}

	sort.SliceStable(hotels, func(i, j int) bool {
		return hotels[i].Rating < hotels[j].Rating
	})
	c.Data["json"] = hotels
	c.ServeJSON()
}

// ReserveRooms ...
// @Title Reserve Rooms
// @Description reserve hotel's rooms
// @Param	body		body 	reqres.BookingReserveRoomsRequest	true		"body for Reservation content"
// @Success 201 {int} reqres.BookingReserveRoomsResponse
// @Failure 403 body is empty
// @router /reserve [post]
func (c *BookingController) ReserveRooms() {
	res := reqres.BookingReserveRoomsResponse{}
	res.SetCode(reqres.Fail)

	defer func() {
		c.Data["json"] = res
		c.ServeJSON()
	}()

	var req reqres.BookingReserveRoomsRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		res.SetCode(reqres.InvalidParams)
		return
	}

	guest, err := models.GetGuestById(req.GuestID)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		res.SetCode(reqres.RecordNotExist)
		return
	}

	reservation := &models.Reservation{
		GuestId:         guest,
		StartDate:       req.StartDate.Time,
		EndDate:         req.EndDate.Time,
		DiscountPercent: req.DiscountPercent,
		TotalPrice:      req.TotalPrice,
		Status:          0,
	}

	o := orm.NewOrm()
	_ = o.Begin()
	// transaction process
	reservationID, err := o.Insert(reservation)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		res.SetCode(reqres.FailedCreate)
		_ = o.Rollback()
		return
	}
	re := models.Reservation{Id: int(reservationID)}
	_ = o.Read(&re)
	for _, roomID := range req.Rooms {
		ro := models.Room{Id: roomID}
		err := o.Read(&ro)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			res.SetCode(reqres.FailedCreate)
			_ = o.Rollback()
			return
		}
		roomReserved := &models.RoomReserved{
			ReservationId: &re,
			RoomId:        &ro,
			Price:         ro.CurrentPrice,
		}
		_, err = o.Insert(roomReserved)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			res.SetCode(reqres.FailedCreate)
			_ = o.Rollback()
			return
		}
	}
	_ = o.Commit()

	c.Ctx.Output.SetStatus(http.StatusCreated)
	res.SetCode(reqres.Success)
	res.Reservation = &re
}
