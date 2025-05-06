package domain

import (
	"time"
)

type UsageType string

const (
	OneTime    UsageType = "one_time"
	MultiUse   UsageType = "multi_use"
	TimeBased  UsageType = "time_based"
)

type DiscountType string

const (
	Percentage DiscountType = "percentage"
	Fixed      DiscountType = "fixed"
)

type Coupon struct {
	ID                    string       `json:"id"`
	CouponCode           string       `json:"coupon_code"`
	ExpiryDate           time.Time    `json:"expiry_date"`
	UsageType            UsageType    `json:"usage_type"`
	ApplicableMedicineIDs []string     `json:"applicable_medicine_ids"`
	ApplicableCategories []string     `json:"applicable_categories"`
	MinOrderValue        float64      `json:"min_order_value"`
	ValidTimeWindow      *TimeWindow  `json:"valid_time_window,omitempty"`
	TermsAndConditions   string       `json:"terms_and_conditions"`
	DiscountType         DiscountType `json:"discount_type"`
	DiscountValue        float64      `json:"discount_value"`
	MaxUsagePerUser      int          `json:"max_usage_per_user"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
}

type TimeWindow struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type CartItem struct {
	ID       string  `json:"id"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type ValidationRequest struct {
	CouponCode  string      `json:"coupon_code"`
	CartItems   []CartItem  `json:"cart_items"`
	OrderTotal  float64     `json:"order_total"`
	Timestamp   time.Time   `json:"timestamp"`
}

type ValidationResponse struct {
	IsValid bool `json:"is_valid"`
	Discount *Discount `json:"discount,omitempty"`
	Message  string    `json:"message"`
	Reason   string    `json:"reason,omitempty"`
}

type Discount struct {
	ItemsDiscount    float64 `json:"items_discount"`
	ChargesDiscount  float64 `json:"charges_discount"`
}

type ApplicableCoupon struct {
	CouponCode    string  `json:"coupon_code"`
	DiscountValue float64 `json:"discount_value"`
} 