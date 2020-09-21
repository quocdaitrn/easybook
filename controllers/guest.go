package controllers

import (
	"crypto/md5"
	"easybook/codecs"
	"easybook/models"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// GuestController operations for Guest
type GuestController struct {
	beego.Controller
}

// URLMapping ...
func (c *GuestController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Guest
// @Param	body		body 	codecs.GuestPostRequest	true		"body for Guest content"
// @Success 201 {int} codecs.GuestPostResponse
// @Failure 403 body is empty
// @router / [post]
func (c *GuestController) Post() {
	res := codecs.GuestPostResponse{}
	res.SetCode(codecs.Fail)

	defer func() {
		c.Data["json"] = res
		c.ServeJSON()
	}()

	var req codecs.GuestPostRequest
	if json.Unmarshal(c.Ctx.Input.RequestBody, &req) != nil {
		c.Ctx.Output.SetStatus(400)
		res.SetCode(codecs.InvalidParams)
		return
	}

	guest := &models.Guest{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Address:   req.Address,
		Detail:    req.Detail,
		Role:      req.Role,
	}

	if req.Password != "" && req.Password == req.ConfirmedPassword {
		hash := md5.Sum([]byte(req.Password))
		guest.Password = hex.EncodeToString(hash[:])
	}

	_, err := models.AddGuest(guest)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		res.SetCode(codecs.FailedCreate)
		return
	}

	c.Ctx.Output.SetStatus(201)
	res.SetCode(codecs.Success)
	res.Guest = guest
}

// GetOne ...
// @Title Get One
// @Description get Guest by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Guest
// @Failure 403 :id is empty
// @router /:id [get]
func (c *GuestController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetGuestById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Guest
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Guest
// @Failure 403
// @router / [get]
func (c *GuestController) GetAll() {
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

	l, err := models.GetAllGuest(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Guest
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Guest	true		"body for Guest content"
// @Success 200 {object} models.Guest
// @Failure 403 :id is not int
// @router /:id [put]
func (c *GuestController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Guest{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateGuestById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Guest
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *GuestController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteGuest(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
