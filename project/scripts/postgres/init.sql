CREATE EXTENSION citext;

CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    email citext UNIQUE NOT NULL,
    password bytea NOT NULL
)