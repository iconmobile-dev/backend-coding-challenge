--

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

-- database schema provision
--

--
-- show only warnings
--

SET client_min_messages TO WARNING;

--
-- create auto-update function and triggers to automatically update
-- updated_at column to NOW()
--

CREATE OR REPLACE FUNCTION set_updated_at_to_now()
    RETURNS TRIGGER AS '
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    ' LANGUAGE 'plpgsql';


--
-- services tables
--

--
-- Getway
--
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id serial,
    email varchar(100) UNIQUE,
    password text NOT NULL,
    firstname text NOT NULL DEFAULT '',
    lastname text NOT NULL DEFAULT '',
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

DROP TRIGGER IF EXISTS users_updated_at ON users;

CREATE TRIGGER users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE set_updated_at_to_now();