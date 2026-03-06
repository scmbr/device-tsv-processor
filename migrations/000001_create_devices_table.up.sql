CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    unit_guid TEXT UNIQUE NOT NULL,
    inv_id TEXT NOT NULL,
    mqtt TEXT,
    status TEXT NOT NULL,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);