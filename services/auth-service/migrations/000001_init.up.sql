-- services/auth-service/migrations/000001_init.up.sql
-- Auth Service Database Schema

-- Permissions table
CREATE TABLE permissions (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    resource VARCHAR(255) NOT NULL,
    action VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_permissions_resource (resource),
    INDEX idx_permissions_action (action),
    INDEX idx_permissions_resource_action (resource, action)
);

-- Roles table
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Role permissions junction table
CREATE TABLE role_permissions (
    role_id VARCHAR(36) NOT NULL,
    permission_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- Tokens table
CREATE TABLE tokens (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type ENUM('access', 'refresh', 'reset', 'verify') NOT NULL,
    value TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    metadata JSON,
    
    INDEX idx_tokens_user_id (user_id),
    INDEX idx_tokens_type (type),
    INDEX idx_tokens_expires_at (expires_at),
    INDEX idx_tokens_value_hash (value(64)),
    INDEX idx_tokens_user_type (user_id, type)
);

-- Sessions table
CREATE TABLE sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address VARCHAR(45),
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_sessions_user_id (user_id),
    INDEX idx_sessions_refresh_token (refresh_token),
    INDEX idx_sessions_expires_at (expires_at),
    INDEX idx_sessions_last_active (last_active)
);

-- User roles junction table (references users from user-service)
CREATE TABLE user_roles (
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by VARCHAR(36),
    
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    INDEX idx_user_roles_user_id (user_id),
    INDEX idx_user_roles_role_id (role_id)
);

-- Insert default permissions
INSERT INTO permissions (id, name, description, resource, action) VALUES
('perm-user-create', 'Create User', 'Create new users', 'user', 'create'),
('perm-user-read', 'Read User', 'View user information', 'user', 'read'),
('perm-user-update', 'Update User', 'Update user information', 'user', 'update'),
('perm-user-delete', 'Delete User', 'Delete users', 'user', 'delete'),
('perm-incident-create', 'Create Incident', 'Create new incidents', 'incident', 'create'),
('perm-incident-read', 'Read Incident', 'View incident information', 'incident', 'read'),
('perm-incident-update', 'Update Incident', 'Update incident information', 'incident', 'update'),
('perm-incident-delete', 'Delete Incident', 'Delete incidents', 'incident', 'delete'),
('perm-incident-assign', 'Assign Incident', 'Assign incidents to users', 'incident', 'assign'),
('perm-admin-access', 'Admin Access', 'Full administrative access', 'system', 'admin');

-- Insert default roles
INSERT INTO roles (id, name, description) VALUES
('role-admin', 'admin', 'System administrator with full access'),
('role-dispatcher', 'dispatcher', 'Dispatcher with incident management access'),
('role-customer', 'customer', 'Customer with limited access');

-- Assign permissions to roles
INSERT INTO role_permissions (role_id, permission_id) VALUES
-- Admin gets all permissions
('role-admin', 'perm-user-create'),
('role-admin', 'perm-user-read'),
('role-admin', 'perm-user-update'),
('role-admin', 'perm-user-delete'),
('role-admin', 'perm-incident-create'),
('role-admin', 'perm-incident-read'),
('role-admin', 'perm-incident-update'),
('role-admin', 'perm-incident-delete'),
('role-admin', 'perm-incident-assign'),
('role-admin', 'perm-admin-access'),

-- Dispatcher gets incident permissions
('role-dispatcher', 'perm-incident-create'),
('role-dispatcher', 'perm-incident-read'),
('role-dispatcher', 'perm-incident-update'),
('role-dispatcher', 'perm-incident-assign'),
('role-dispatcher', 'perm-user-read'),

-- Customer gets basic read permissions
('role-customer', 'perm-incident-read'),
('role-customer', 'perm-user-read');