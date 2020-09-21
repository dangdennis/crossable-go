--
-- PLEASE NOTE - this is destructive action. uncomment and execute only if you are really sure about dropping tables.
--

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
COMMENT ON SCHEMA public IS 'standard public schema';