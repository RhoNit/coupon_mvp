package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/RhoNit/coupon_mvp/internal/domain"
	"github.com/RhoNit/coupon_mvp/internal/service"
)

type CouponHandler struct {
	service *service.CouponService
}

func NewCouponHandler(service *service.CouponService) *CouponHandler {
	return &CouponHandler{service: service}
}

// @Summary Get applicable coupons
// @Description Get all applicable coupons for the given cart
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body domain.ValidationRequest true "Cart details"
// @Success 200 {array} domain.ApplicableCoupon
// @Router /coupons/applicable [get]
func (h *CouponHandler) GetApplicableCoupons(c echo.Context) error {
	var req domain.ValidationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	coupons, err := h.service.GetApplicableCoupons(c.Request().Context(), req.CartItems, req.OrderTotal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get applicable coupons",
		})
	}

	return c.JSON(http.StatusOK, coupons)
}

// @Summary Validate coupon
// @Description Validate a coupon code against cart items
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body domain.ValidationRequest true "Validation request"
// @Success 200 {object} domain.ValidationResponse
// @Router /coupons/validate [post]
func (h *CouponHandler) ValidateCoupon(c echo.Context) error {
	var req domain.ValidationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	response, err := h.service.ValidateCoupon(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to validate coupon",
		})
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary Create coupon
// @Description Create a new coupon
// @Tags coupons
// @Accept json
// @Produce json
// @Param coupon body domain.Coupon true "Coupon details"
// @Success 201 {object} domain.Coupon
// @Router /coupons [post]
func (h *CouponHandler) CreateCoupon(c echo.Context) error {
	var coupon domain.Coupon
	if err := c.Bind(&coupon); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Set default values
	coupon.ID = uuid.New().String()
	coupon.CreatedAt = time.Now()
	coupon.UpdatedAt = time.Now()

	if err := h.service.CreateCoupon(c.Request().Context(), &coupon); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create coupon",
		})
	}

	return c.JSON(http.StatusCreated, coupon)
}
