-- services/user-service/migrations/000001_init.down.sql
-- Rollback User Service Database Schema

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS password_reset_attempts;
DROP TABLE IF EXISTS user_activity_log;
DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS user_contacts;
DROP TABLE IF EXISTS users;