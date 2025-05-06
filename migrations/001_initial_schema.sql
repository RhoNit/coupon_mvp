-- Create coupons table
CREATE TABLE IF NOT EXISTS coupons (
    id VARCHAR(36) PRIMARY KEY,
    coupon_code VARCHAR(50) UNIQUE NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    usage_type VARCHAR(20) NOT NULL,
    applicable_medicine_ids TEXT[],
    applicable_categories TEXT[],
    min_order_value DECIMAL(10,2) NOT NULL,
    valid_time_window_start TIMESTAMP,
    valid_time_window_end TIMESTAMP,
    terms_and_conditions TEXT,
    discount_type VARCHAR(20) NOT NULL,
    discount_value DECIMAL(10,2) NOT NULL,
    max_usage_per_user INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create coupon_usage table
CREATE TABLE IF NOT EXISTS coupon_usage (
    id SERIAL PRIMARY KEY,
    coupon_code VARCHAR(50) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    used_at TIMESTAMP NOT NULL,
    FOREIGN KEY (coupon_code) REFERENCES coupons(coupon_code),
    UNIQUE(coupon_code, user_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_coupons_code ON coupons(coupon_code);
CREATE INDEX IF NOT EXISTS idx_coupons_expiry ON coupons(expiry_date);
CREATE INDEX IF NOT EXISTS idx_coupon_usage_code ON coupon_usage(coupon_code);
CREATE INDEX IF NOT EXISTS idx_coupon_usage_user ON coupon_usage(user_id); 