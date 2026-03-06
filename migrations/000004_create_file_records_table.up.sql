CREATE TABLE file_records (
    id SERIAL PRIMARY KEY,
    filename TEXT UNIQUE NOT NULL,
    full_path TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP,
    attempts INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL,
    error_message TEXT,
    updated_at TIMESTAMP
);