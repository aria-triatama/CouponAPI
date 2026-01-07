package main

import (
	"CouponAPI/helpers"
	"CouponAPI/route"
)

func main() {
	e := route.Init()
	e.Logger.Fatal(e.Start(":" + helpers.GetEnv("PORT")))
	route.Mdb.Disconnect()
}
