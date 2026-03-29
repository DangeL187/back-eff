CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    service_name TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL DEFAULT 'infinity'::date,
    price_per_day BIGINT NOT NULL
);

ALTER TABLE subscriptions
    ADD COLUMN period tsrange GENERATED ALWAYS AS (tsrange(start_date, end_date, '[]')) STORED;

CREATE INDEX idx_subscriptions_period ON subscriptions USING GIST (period);
CREATE INDEX idx_subscriptions_user_service ON subscriptions(user_id, service_name);