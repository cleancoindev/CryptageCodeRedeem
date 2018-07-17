package endpoint

import (
	"strconv"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/DecenterApps/CryptageCodeRedeem/app"
	"github.com/DecenterApps/CryptageCodeRedeem/model"
)

type (
	// couponService specifies the interface for the coupon service needed by couponResource.
	couponService interface {
		Get(rs app.RequestScope, id int) (*model.Coupon, error)
		GetByToken(rs app.RequestScope, token int) (*model.Coupon, error)
		Query(rs app.RequestScope, offset, limit int) ([]model.Coupon, error)
		Count(rs app.RequestScope) (int, error)
		Create(rs app.RequestScope, model *model.Coupon) (*model.Coupon, error)
		Update(rs app.RequestScope, id int, model *model.Coupon) (*model.Coupon, error)
		Delete(rs app.RequestScope, id int) (*model.Coupon, error)
	}

	couponResource struct {
		service couponService
	}
)

func ServeCouponResource(rg *routing.RouteGroup, service couponService) {
	r := &couponResource{service}
	rg.Get("/freeCard/<token>", r.get)
	rg.Post("/freeCard/<token>", r.update)
}

func (r *couponResource) get(c *routing.Context) error {
	token, err := strconv.Atoi(c.Param("token"))
	if err != nil {
		return err
	}

	coupon, err := r.service.GetByToken(app.GetRequestScope(c), token)
	if err != nil {
		return err
	}

	if coupon.User != nil {
		return c.Write("Coupon already used")
	}

	return c.Write(coupon)
}

func (r *couponResource) update(c *routing.Context) error {
	token, err := strconv.Atoi(c.Param("token"))
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)

	coupon, err := r.service.GetByToken(rs, token)

	if coupon.User != nil {
		return c.Write("Coupon already used")
	}
	if err != nil {
		return err
	}

	if err := c.Read(coupon); err != nil {
		return err
	}

	response, err := r.service.Update(rs, coupon.Id, coupon)
	if err != nil {
		return err
	}

	return c.Write(response)
}