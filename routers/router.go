// @APIVersion 1.0.0
// @Title Hotel Booking API
// @Description Autogenerate API documents
// @Contact daitq.cntt@gmail.com
package routers

import (
	"easybook/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/rpc",
			beego.NSNamespace("/hotels",
				beego.NSInclude(
					&controllers.BookingController{},
				),
			),
			beego.NSNamespace("/rooms",
				beego.NSInclude(
					&controllers.BookingController{},
				),
			),
		),

		beego.NSNamespace("/guests",
			beego.NSInclude(
				&controllers.GuestController{},
			),
		),

		beego.NSNamespace("/hotels",
			beego.NSInclude(
				&controllers.HotelController{},
			),
		),

		beego.NSNamespace("/notifications",
			beego.NSInclude(
				&controllers.NotificationController{},
			),
		),

		beego.NSNamespace("/reservations",
			beego.NSInclude(
				&controllers.ReservationController{},
			),
		),

		beego.NSNamespace("/rooms",
			beego.NSInclude(
				&controllers.RoomController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
