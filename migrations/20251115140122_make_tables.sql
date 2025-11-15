CREATE TABLE teams (
  id SERIAL PRIMARY KEY,
  team_name TEXT UNIQUE NOT NULL
);

CREATE TABLE users (
  id TEXT PRIMARY KEY,        
  username TEXT NOT NULL,
  team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
  is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE pull_requests (
  id TEXT PRIMARY KEY,         
  name TEXT NOT NULL,
  author_id TEXT NOT NULL REFERENCES users(id),
  status TEXT NOT NULL CHECK (status IN ('OPEN','MERGED')) DEFAULT 'OPEN',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  merged_at TIMESTAMPTZ NULL
);
