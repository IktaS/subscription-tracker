CREATE TABLE IF NOT EXISTS subscription(
    id TEXT not null primary key,
    user_id TEXT not null,
    title TEXT,
    description TEXT,
    payment_method TEXT,
    amount_currency VARCHAR(3),
    amount_value DOUBLE,
    last_paid DATE,
    duration_value INT,
    duration_unit TEXT
);
CREATE TABLE IF NOT EXISTS configuration(log_channel TEXT);