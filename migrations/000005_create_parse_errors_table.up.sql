CREATE TABLE parse_error_model (
    id SERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    file_id INTEGER NOT NULL REFERENCES file_record(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);