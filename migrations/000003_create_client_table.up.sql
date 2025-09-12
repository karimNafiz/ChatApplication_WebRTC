-- creating the clients table, 
CREATE TABLE app.clients(
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);