package controllers

import (
	"net/http"

	"CouponAPI/database"
	"CouponAPI/models"

	"github.com/labstack/echo/v4"
)

var Mdb, _ = database.MongoDB()

func CreateCoupon(c echo.Context) error {

	return c.NoContent(http.StatusCreated)
}

func ClaimCoupon(c echo.Context) error {

	return c.NoContent(http.StatusOK)
}

func GetCouponDetails(c echo.Context) error {
	couponName := c.Param("name")
	if couponName == "" {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, models.CouponDetails{Name: couponName})
}
