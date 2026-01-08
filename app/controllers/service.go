package controllers

import (
	"net/http"

	"CouponAPI/database"
	"CouponAPI/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var Mdb, _ = database.MongoDB()
var validate = validator.New(validator.WithRequiredStructEnabled())

func CreateCoupon(c echo.Context) error {
	var coupon models.CouponsRequest

	if err := c.Bind(&coupon); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := validate.Struct(coupon); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	_, err := database.AddCoupon(Mdb, models.Coupons{
		Name:            coupon.Name,
		Amount:          coupon.Amount,
		RemainingAmount: coupon.Amount,
	})
	if err != nil {
		if err.Error() == "coupon already exists" {
			return c.NoContent(http.StatusConflict)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func ClaimCoupon(c echo.Context) error {
	var claim models.Claims

	if err := c.Bind(&claim); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	if err := validate.Struct(claim); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if database.CheckClaim(Mdb, claim) {
		return c.NoContent(http.StatusConflict)
	}

	if !database.CheckStock(Mdb, claim) {
		return c.NoContent(http.StatusBadRequest)
	}

	if !database.ReserveGate(Mdb, claim.CouponName) {
		return c.NoContent(http.StatusInternalServerError)
	}

	if err := database.AddClaim(Mdb, claim); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if err := database.UpdateStock(Mdb, claim); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	_ = database.ClearGate(Mdb, claim.CouponName)

	return c.NoContent(http.StatusOK)
}

func GetCouponDetails(c echo.Context) error {
	couponName := c.Param("name")
	if couponName == "" {
		return c.NoContent(http.StatusNotFound)
	}

	detail, err := database.GetCouponDetails(Mdb, couponName)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, detail)
}
