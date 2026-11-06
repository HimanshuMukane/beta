-- Initialize roles and databases for Supabase / PostgreSQL setup
-- Primary superuser / application credentials: username=especa, password=especa

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'especa') THEN
        CREATE ROLE especa WITH LOGIN SUPERUSER CREATEDB CREATEROLE REPLICATION PASSWORD 'especa';
    ELSE
        ALTER ROLE especa WITH LOGIN SUPERUSER CREATEDB CREATEROLE REPLICATION PASSWORD 'especa';
    END IF;
END
$$;

