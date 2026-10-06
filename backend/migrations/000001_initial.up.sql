CREATE TABLE customers (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(254) NOT NULL UNIQUE,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB;

CREATE TABLE merchants (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB;

CREATE TABLE orders (
  id CHAR(36) PRIMARY KEY,
  reference VARCHAR(32) NOT NULL UNIQUE,
  customer_id CHAR(36) NOT NULL,
  merchant_id CHAR(36) NOT NULL,
  amount BIGINT NOT NULL,
  currency CHAR(3) NOT NULL,
  status ENUM('pending','active','completed','cancelled','refunded') NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  CONSTRAINT chk_order_amount CHECK (amount > 0),
  CONSTRAINT chk_order_currency CHECK (currency IN ('USD','CLP')),
  CONSTRAINT fk_order_customer FOREIGN KEY (customer_id) REFERENCES customers(id),
  CONSTRAINT fk_order_merchant FOREIGN KEY (merchant_id) REFERENCES merchants(id),
  INDEX idx_orders_customer_created (customer_id, created_at DESC),
  INDEX idx_orders_merchant_created (merchant_id, created_at DESC)
) ENGINE=InnoDB;

CREATE TABLE payment_plans (
  id CHAR(36) PRIMARY KEY,
  order_id CHAR(36) NOT NULL UNIQUE,
  installment_count TINYINT UNSIGNED NOT NULL,
  status ENUM('active','completed','refunded') NOT NULL DEFAULT 'active',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  CONSTRAINT chk_installment_count CHECK (installment_count IN (2,3,4)),
  CONSTRAINT fk_plan_order FOREIGN KEY (order_id) REFERENCES orders(id)
) ENGINE=InnoDB;

CREATE TABLE installments (
  id CHAR(36) PRIMARY KEY,
  payment_plan_id CHAR(36) NOT NULL,
  sequence_number TINYINT UNSIGNED NOT NULL,
  amount BIGINT NOT NULL,
  status ENUM('pending','paid','failed','refunded') NOT NULL DEFAULT 'pending',
  due_at TIMESTAMP(6) NOT NULL,
  paid_at TIMESTAMP(6) NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  CONSTRAINT chk_installment_amount CHECK (amount > 0),
  CONSTRAINT fk_installment_plan FOREIGN KEY (payment_plan_id) REFERENCES payment_plans(id),
  UNIQUE KEY uq_installment_sequence (payment_plan_id, sequence_number),
  INDEX idx_installments_due (status, due_at)
) ENGINE=InnoDB;

CREATE TABLE payments (
  id CHAR(36) PRIMARY KEY,
  installment_id CHAR(36) NOT NULL,
  amount BIGINT NOT NULL,
  currency CHAR(3) NOT NULL,
  status ENUM('pending','succeeded','failed','refunded') NOT NULL,
  failure_reason VARCHAR(255) NOT NULL DEFAULT '',
  refunded_amount BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  CONSTRAINT chk_payment_amount CHECK (amount > 0),
  CONSTRAINT chk_payment_refunded CHECK (refunded_amount >= 0 AND refunded_amount <= amount),
  CONSTRAINT fk_payment_installment FOREIGN KEY (installment_id) REFERENCES installments(id),
  INDEX idx_payments_installment_created (installment_id, created_at DESC)
) ENGINE=InnoDB;

CREATE TABLE refunds (
  id CHAR(36) PRIMARY KEY,
  payment_id CHAR(36) NOT NULL,
  amount BIGINT NOT NULL,
  reason VARCHAR(255) NOT NULL,
  status ENUM('succeeded','failed') NOT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  CONSTRAINT chk_refund_amount CHECK (amount > 0),
  CONSTRAINT fk_refund_payment FOREIGN KEY (payment_id) REFERENCES payments(id),
  INDEX idx_refunds_payment_created (payment_id, created_at DESC)
) ENGINE=InnoDB;

CREATE TABLE transaction_events (
  id CHAR(36) PRIMARY KEY,
  customer_id CHAR(36) NOT NULL,
  order_id CHAR(36) NOT NULL,
  event_type ENUM('order_created','plan_created','payment_succeeded','payment_failed','refund_succeeded') NOT NULL,
  status VARCHAR(32) NOT NULL,
  amount BIGINT NOT NULL,
  currency CHAR(3) NOT NULL,
  reference_id CHAR(36) NOT NULL,
  metadata JSON NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  CONSTRAINT fk_event_customer FOREIGN KEY (customer_id) REFERENCES customers(id),
  CONSTRAINT fk_event_order FOREIGN KEY (order_id) REFERENCES orders(id),
  INDEX idx_events_customer_created (customer_id, created_at DESC),
  INDEX idx_events_order_created (order_id, created_at DESC)
) ENGINE=InnoDB;

CREATE TABLE idempotency_keys (
  idempotency_key VARCHAR(128) PRIMARY KEY,
  request_hash CHAR(64) NOT NULL,
  resource_type VARCHAR(32) NOT NULL,
  resource_id CHAR(36) NULL,
  response_code SMALLINT UNSIGNED NULL,
  state ENUM('processing','completed') NOT NULL DEFAULT 'processing',
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  expires_at TIMESTAMP(6) NOT NULL,
  INDEX idx_idempotency_expiry (expires_at)
) ENGINE=InnoDB;

