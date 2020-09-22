package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error400() {
	c.Ctx.Output.SetStatus(http.StatusBadRequest)
	c.Ctx.Output.Body([]byte("Bad request."))
}

func (c *ErrorController) Error404() {
	c.Ctx.Output.SetStatus(http.StatusNotFound)
	c.Ctx.Output.Body([]byte("Not found."))
}

func (c *ErrorController) Error500() {
	c.Ctx.Output.SetStatus(http.StatusInternalServerError)
	c.Ctx.Output.Body([]byte("Internal server error."))
}

func (c *ErrorController) ErrorDb() {
	c.Ctx.Output.SetStatus(http.StatusInternalServerError)
	c.Ctx.Output.Body([]byte("Database is now down."))
}

func (c *ErrorController) ErrorDashboard401() {
	c.Ctx.Output.SetStatus(http.StatusUnauthorized)
	c.Ctx.Output.Body([]byte("User is not authorized."))
}
