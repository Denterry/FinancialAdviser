-- Drop tables
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS plans;

-- Drop enum types
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS plan_type; 