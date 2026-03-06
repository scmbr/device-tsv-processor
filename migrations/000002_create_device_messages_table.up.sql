CREATE TABLE device_messages (
    id SERIAL PRIMARY KEY,
    device_id INTEGER NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    inv_id TEXT NOT NULL,
    msg_id TEXT UNIQUE NOT NULL,
    text TEXT,
    context TEXT,
    class TEXT,
    level INTEGER,
    area TEXT,
    addr TEXT,
    block TEXT,
    type TEXT,
    bit INTEGER,
    invert_bit BOOLEAN,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);