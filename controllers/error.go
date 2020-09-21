package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error400() {
	c.Ctx.Output.SetStatus(400)
	c.Ctx.Output.Body([]byte("Bad request."))
}

func (c *ErrorController) Error404() {
	c.Ctx.Output.SetStatus(404)
	c.Ctx.Output.Body([]byte("Page not found."))
}

func (c *ErrorController) Error500() {
	c.Ctx.Output.SetStatus(500)
	c.Ctx.Output.Body([]byte("Internal server error."))
}

func (c *ErrorController) ErrorDb() {
	c.Ctx.Output.SetStatus(500)
	c.Ctx.Output.Body([]byte("Database is now down."))
}

func (c *ErrorController) ErrorDevice401() {
	c.Ctx.Output.SetStatus(401)
	c.Ctx.Output.Body([]byte("Device not authorized."))
}

func (c *ErrorController) ErrorDashboard401() {
	c.Ctx.Output.SetStatus(401)
	c.Ctx.Output.Body([]byte("User not authorized."))
}
