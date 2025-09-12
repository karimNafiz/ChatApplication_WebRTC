-- creating our app user 
CREATE ROLE app_user
    WITH LOGIN PASSWORD 'pa55word'
    NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT NOREPLICATION;

-- grant app_user access to chat_application
GRANT CONNECT ON DATABASE chat_application TO app_user;

-- let the user use the default schema
GRANT USAGE ON SCHEMA public TO app_user;

-- give rights to all the current tables, right now there is nothing but its fine
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;

-- grant access on all sequences
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO app_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public 
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO app_user;