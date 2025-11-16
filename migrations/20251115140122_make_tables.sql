-- +goose Up
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    team_name TEXT UNIQUE NOT NULL
);

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE pull_requests (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

    status TEXT NOT NULL CHECK (status IN ('OPEN','MERGED')) DEFAULT 'OPEN',

    assigned_reviewers TEXT[] NOT NULL DEFAULT '{}',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    merged_at TIMESTAMPTZ NULL
);


-- +goose Down
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;