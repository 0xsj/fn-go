-- services/auth-service/migrations/000001_init.down.sql
-- Rollback Auth Service Database Schema

-- Drop tables in reverse order of creation (due to foreign key constraints)
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;