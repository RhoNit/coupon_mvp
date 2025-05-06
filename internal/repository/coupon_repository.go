package repository

import (
	"context"
	"errors"
	"time"

	"github.com/RhoNit/coupon_mvp/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CouponRepository interface {
	Create(ctx context.Context, coupon *domain.Coupon) error
	GetByCode(ctx context.Context, code string) (*domain.Coupon, error)
	GetApplicableCoupons(ctx context.Context, cartItems []domain.CartItem, orderTotal float64) ([]domain.ApplicableCoupon, error)
	UpdateUsage(ctx context.Context, couponCode string, userId string) error
	GetUsageCount(ctx context.Context, couponCode string, userId string) (int, error)
}

type couponRepository struct {
	pool *pgxpool.Pool
}

func NewCouponRepository(pool *pgxpool.Pool) CouponRepository {
	return &couponRepository{pool: pool}
}

func (r *couponRepository) Create(ctx context.Context, coupon *domain.Coupon) error {
	query := `
		INSERT INTO coupons (
			id, coupon_code, expiry_date, usage_type, applicable_medicine_ids,
			applicable_categories, min_order_value, valid_time_window_start,
			valid_time_window_end, terms_and_conditions, discount_type,
			discount_value, max_usage_per_user, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	_, err := r.pool.Exec(ctx, query,
		coupon.ID,
		coupon.CouponCode,
		coupon.ExpiryDate,
		coupon.UsageType,
		coupon.ApplicableMedicineIDs,
		coupon.ApplicableCategories,
		coupon.MinOrderValue,
		coupon.ValidTimeWindow.StartTime,
		coupon.ValidTimeWindow.EndTime,
		coupon.TermsAndConditions,
		coupon.DiscountType,
		coupon.DiscountValue,
		coupon.MaxUsagePerUser,
		time.Now(),
		time.Now(),
	)

	return err
}

func (r *couponRepository) GetByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	query := `
		SELECT id, coupon_code, expiry_date, usage_type, applicable_medicine_ids,
			applicable_categories, min_order_value, valid_time_window_start,
			valid_time_window_end, terms_and_conditions, discount_type,
			discount_value, max_usage_per_user, created_at, updated_at
		FROM coupons
		WHERE coupon_code = $1
	`

	var coupon domain.Coupon
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&coupon.ID,
		&coupon.CouponCode,
		&coupon.ExpiryDate,
		&coupon.UsageType,
		&coupon.ApplicableMedicineIDs,
		&coupon.ApplicableCategories,
		&coupon.MinOrderValue,
		&coupon.ValidTimeWindow.StartTime,
		&coupon.ValidTimeWindow.EndTime,
		&coupon.TermsAndConditions,
		&coupon.DiscountType,
		&coupon.DiscountValue,
		&coupon.MaxUsagePerUser,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &coupon, nil
}

func (r *couponRepository) GetApplicableCoupons(ctx context.Context, cartItems []domain.CartItem, orderTotal float64) ([]domain.ApplicableCoupon, error) {
	query := `
		SELECT coupon_code, discount_value
		FROM coupons
		WHERE expiry_date > NOW()
		AND min_order_value <= $1
		AND (
			valid_time_window_start IS NULL 
			OR (NOW() BETWEEN valid_time_window_start AND valid_time_window_end)
		)
	`

	rows, err := r.pool.Query(ctx, query, orderTotal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []domain.ApplicableCoupon
	for rows.Next() {
		var coupon domain.ApplicableCoupon
		if err := rows.Scan(&coupon.CouponCode, &coupon.DiscountValue); err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}

	return coupons, nil
}

func (r *couponRepository) UpdateUsage(ctx context.Context, couponCode string, userId string) error {
	query := `
		INSERT INTO coupon_usage (coupon_code, user_id, used_at)
		VALUES ($1, $2, NOW())
	`

	_, err := r.pool.Exec(ctx, query, couponCode, userId)
	return err
}

func (r *couponRepository) GetUsageCount(ctx context.Context, couponCode string, userId string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM coupon_usage
		WHERE coupon_code = $1 AND user_id = $2
	`

	var count int
	err := r.pool.QueryRow(ctx, query, couponCode, userId).Scan(&count)
	return count, err
}
