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
		),

		beego.NSNamespace("/agreements",
			beego.NSInclude(
				&controllers.AgreementController{},
			),
		),

		beego.NSNamespace("/cities",
			beego.NSInclude(
				&controllers.CityController{},
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

		beego.NSNamespace("/invoices",
			beego.NSInclude(
				&controllers.InvoiceController{},
			),
		),

		beego.NSNamespace("/notifications",
			beego.NSInclude(
				&controllers.NotificationController{},
			),
		),

		beego.NSNamespace("/penalty_rules",
			beego.NSInclude(
				&controllers.PenaltyRuleController{},
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

		beego.NSNamespace("/room_facilitates",
			beego.NSInclude(
				&controllers.RoomFacilitateController{},
			),
		),

		beego.NSNamespace("/room_reserveds",
			beego.NSInclude(
				&controllers.RoomReservedController{},
			),
		),

		beego.NSNamespace("/service_level",
			beego.NSInclude(
				&controllers.ServiceLevelController{},
			),
		),

		beego.NSNamespace("/stay_tracking",
			beego.NSInclude(
				&controllers.StayTrackingController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
