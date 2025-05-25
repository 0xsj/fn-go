-- services/user-service/migrations/000001_init.up.sql
-- User Service Database Schema

-- Users table
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone VARCHAR(20),
    role ENUM('admin', 'customer', 'dispatcher') NOT NULL DEFAULT 'customer',
    status ENUM('active', 'inactive', 'suspended', 'pending') NOT NULL DEFAULT 'pending',
    last_login_at TIMESTAMP NULL,
    failed_logins INT DEFAULT 0,
    email_verified BOOLEAN DEFAULT FALSE,
    preferences JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_users_username (username),
    INDEX idx_users_email (email),
    INDEX idx_users_role (role),
    INDEX idx_users_status (status),
    INDEX idx_users_email_verified (email_verified),
    INDEX idx_users_deleted_at (deleted_at),
    INDEX idx_users_created_at (created_at)
);

-- User contacts table (for additional contact methods)
CREATE TABLE user_contacts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type ENUM('email', 'phone', 'mobile', 'work_phone', 'emergency') NOT NULL,
    value VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_contacts_user_id (user_id),
    INDEX idx_user_contacts_type (type),
    INDEX idx_user_contacts_primary (is_primary),
    UNIQUE KEY uk_user_contacts_primary (user_id, type, is_primary)
);

-- User preferences table (for complex preferences)
CREATE TABLE user_preferences (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE,
    theme VARCHAR(50) DEFAULT 'light',
    notifications_enabled BOOLEAN DEFAULT TRUE,
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    email_notifications BOOLEAN DEFAULT TRUE,
    sms_notifications BOOLEAN DEFAULT FALSE,
    push_notifications BOOLEAN DEFAULT TRUE,
    incident_notifications BOOLEAN DEFAULT TRUE,
    chat_notifications BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- User activity log
CREATE TABLE user_activity_log (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    action VARCHAR(255) NOT NULL,
    resource VARCHAR(255),
    resource_id VARCHAR(36),
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_activity_user_id (user_id),
    INDEX idx_activity_action (action),
    INDEX idx_activity_resource (resource),
    INDEX idx_activity_created_at (created_at)
);

-- Password reset attempts
CREATE TABLE password_reset_attempts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45),
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_password_reset_user_id (user_id),
    INDEX idx_password_reset_token (token_hash),
    INDEX idx_password_reset_expires (expires_at)
);

-- Insert default admin user (password: admin123 - change in production!)
INSERT INTO users (
    id, 
    username, 
    email, 
    password_hash, 
    first_name, 
    last_name, 
    role, 
    status, 
    email_verified
) VALUES (
    'user-admin-default',
    'admin',
    'admin@fn-go.local',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- bcrypt hash of 'admin123'
    'System',
    'Administrator',
    'admin',
    'active',
    TRUE
);

-- Insert preferences for admin user
INSERT INTO user_preferences (
    id,
    user_id,
    theme,
    notifications_enabled,
    language,
    timezone
) VALUES (
    'pref-admin-default',
    'user-admin-default',
    'dark',
    TRUE,
    'en',
    'UTC'
);