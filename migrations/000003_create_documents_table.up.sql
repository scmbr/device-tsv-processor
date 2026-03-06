CREATE TABLE document (
    id SERIAL PRIMARY KEY,
    unit_guid TEXT NOT NULL,
    file_path TEXT,
    file_type TEXT NOT NULL,
    status TEXT NOT NULL,
    attempts INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);