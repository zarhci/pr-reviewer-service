-- +goose Up
CREATE TABLE teams (
    team_name TEXT PRIMARY KEY
);

CREATE TABLE users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    team_name TEXT REFERENCES teams(team_name) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE pull_requests (
    pull_request_id TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(user_id),
    status TEXT NOT NULL,
    created_at TIMESTAMP,
    merged_at TIMESTAMP
);

CREATE TABLE pull_request_reviewers (
    pull_request_id TEXT REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id TEXT REFERENCES users(user_id),
    PRIMARY KEY (pull_request_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS pull_request_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
