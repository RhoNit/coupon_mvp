package service

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/RhoNit/coupon_mvp/internal/domain"
	"github.com/RhoNit/coupon_mvp/internal/repository"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCouponNotFound    = errors.New("coupon not found")
	ErrCouponExpired     = errors.New("coupon has expired")
	ErrInvalidUsageType  = errors.New("invalid usage type")
	ErrMaxUsageExceeded  = errors.New("maximum usage limit exceeded")
	ErrMinOrderNotMet    = errors.New("minimum order value not met")
	ErrInvalidTimeWindow = errors.New("invalid time window")
	ErrInvalidMedicine   = errors.New("invalid medicine for coupon")
	ErrInvalidCategory   = errors.New("invalid category for coupon")
)

type CouponService struct {
	repo  repository.CouponRepository
	redis *redis.Client
	mu    sync.RWMutex
}

func NewCouponService(repo repository.CouponRepository, redisClient *redis.Client) *CouponService {
	return &CouponService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *CouponService) ValidateCoupon(ctx context.Context, req *domain.ValidationRequest) (*domain.ValidationResponse, error) {
	// try to get coupon from redis first
	cacheKey := "coupon:" + req.CouponCode
	cachedData, err := s.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var coupon domain.Coupon
		if err := json.Unmarshal(cachedData, &coupon); err == nil {
			return s.validateCouponRules(ctx, &coupon, req)
		}
	}

	// if not present in redis, get from database
	coupon, err := s.repo.GetByCode(ctx, req.CouponCode)
	if err != nil {
		return nil, err
	}
	if coupon == nil {
		return &domain.ValidationResponse{
			IsValid: false,
			Reason:  ErrCouponNotFound.Error(),
		}, nil
	}

	// cache the coupon in redis
	couponData, err := json.Marshal(coupon)
	if err == nil {
		s.redis.Set(ctx, cacheKey, couponData, 5*time.Minute)
	}

	return s.validateCouponRules(ctx, coupon, req)
}

func (s *CouponService) validateCouponRules(ctx context.Context, coupon *domain.Coupon, req *domain.ValidationRequest) (*domain.ValidationResponse, error) {
	// check the expiry
	if time.Now().After(coupon.ExpiryDate) {
		return &domain.ValidationResponse{
			IsValid: false,
			Reason:  ErrCouponExpired.Error(),
		}, nil
	}

	// check minimum order value
	if req.OrderTotal < coupon.MinOrderValue {
		return &domain.ValidationResponse{
			IsValid: false,
			Reason:  ErrMinOrderNotMet.Error(),
		}, nil
	}

	// check time window if applicable
	if coupon.ValidTimeWindow != nil {
		now := time.Now()
		if now.Before(coupon.ValidTimeWindow.StartTime) || now.After(coupon.ValidTimeWindow.EndTime) {
			return &domain.ValidationResponse{
				IsValid: false,
				Reason:  ErrInvalidTimeWindow.Error(),
			}, nil
		}
	}

	// check applicable medicines and categories
	validItems := false
	for _, item := range req.CartItems {
		// check medicine IDs
		for _, medID := range coupon.ApplicableMedicineIDs {
			if item.ID == medID {
				validItems = true
				break
			}
		}
		// check categories
		for _, category := range coupon.ApplicableCategories {
			if item.Category == category {
				validItems = true
				break
			}
		}
	}

	if !validItems {
		return &domain.ValidationResponse{
			IsValid: false,
			Reason:  "No applicable items in cart",
		}, nil
	}

	// calculate discount
	discount := &domain.Discount{}
	if coupon.DiscountType == domain.Percentage {
		discount.ItemsDiscount = req.OrderTotal * (coupon.DiscountValue / 100)
	} else {
		discount.ItemsDiscount = coupon.DiscountValue
	}

	return &domain.ValidationResponse{
		IsValid:  true,
		Discount: discount,
		Message:  "Coupon applied successfully",
	}, nil
}

func (s *CouponService) GetApplicableCoupons(ctx context.Context, cartItems []domain.CartItem, orderTotal float64) ([]domain.ApplicableCoupon, error) {
	return s.repo.GetApplicableCoupons(ctx, cartItems, orderTotal)
}

func (s *CouponService) CreateCoupon(ctx context.Context, coupon *domain.Coupon) error {
	return s.repo.Create(ctx, coupon)
}
