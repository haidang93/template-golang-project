CREATE TABLE IF NOT EXISTS file (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    url TEXT NOT NULL,
    path TEXT NOT NULL,
    type TEXT NOT NULL,
    file_name TEXT NOT NULL,
    extension citext NOT NULL,
    size BIGINT,
    create_date TIMESTAMPTZ DEFAULT NOW(),
    update_date TIMESTAMPTZ
);