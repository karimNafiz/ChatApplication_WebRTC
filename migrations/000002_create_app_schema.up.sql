-- create the app schema 
CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION postgres;

-- let our app_user to see/use the schema
GRANT USAGE ON SCHEMA app TO app_user;

-- allow data operations for app user
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA app to app_user;

-- allow operations on sequences
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA app to app_user;

-- ensure that same permissions are applied to tables/sequences in this schema
ALTER DEFAULT PRIVILEGES IN SCHEMA app
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA app
  GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO app_user;