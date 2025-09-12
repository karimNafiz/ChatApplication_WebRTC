-- the cascade command makes sure that, everything inside the schema is also deleted
-- if there are tables, functions, sequences or other objects, then without CASCADE this command won't work
DROP SCHEMA IF EXISTS app CASCADE;