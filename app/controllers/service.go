package controllers

import (
	"errors"
	"net/http"

	"CouponAPI/database"
	"CouponAPI/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
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

	client, _ := Mdb.GetClient()

	session, err := client.StartSession()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	defer session.EndSession(*Mdb.GetCtx())

	err = mongo.WithSession(*Mdb.GetCtx(), session, func(s mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if !database.CheckStock(Mdb, s, claim) {
			session.AbortTransaction(s)
			return errors.New("no stock")
		}

		if err := database.AddClaim(Mdb, s, claim); err != nil {
			session.AbortTransaction(s)
			return err
		}

		if err := database.UpdateStock(Mdb, s, claim); err != nil {
			session.AbortTransaction(s)
			return err
		}

		return session.CommitTransaction(s)
	})

	if err != nil {
		if err.Error() == "no stock" {
			return c.NoContent(http.StatusBadRequest)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

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
