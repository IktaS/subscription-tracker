CREATE TABLE IF NOT EXISTS user(
    id TEXT not null primary key,
    payday_time TEXT
);
CREATE TABLE IF NOT EXISTS subscription(
    id TEXT not null primary key,
    user_id TEXT not null,
    title TEXT,
    payment_method TEXT,
    amount_currency VARCHAR(3),
    amount_value DOUBLE,
    last_paid DATETIME,
    duration_value INT,
    duration_unit TEXT,
    FOREIGN KEY(user_id) REFERENCES user(id)
);
CREATE TABLE IF NOT EXISTS configuration(log_channel TEXT);