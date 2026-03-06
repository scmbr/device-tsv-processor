CREATE TABLE file_record (
    id SERIAL PRIMARY KEY,
    filename TEXT UNIQUE NOT NULL,
    full_path TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP,
    status TEXT NOT NULL,
    error_message TEXT,
    updated_at TIMESTAMP
);