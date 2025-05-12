// pkg/models/auth.go
package models

import "time"

type TokenType string

const (
    TokenTypeAccess  TokenType = "access"
    TokenTypeRefresh TokenType = "refresh"
    TokenTypeReset   TokenType = "reset"
    TokenTypeVerify  TokenType = "verify"
)

type Token struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Type      TokenType `json:"type"`
    Value     string    `json:"value"`
    ExpiresAt time.Time `json:"expires_at"`
    RevokedAt *time.Time `json:"revoked_at,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    Metadata  map[string]any `json:"metadata,omitempty"`
}

type TokenClaims struct {
    UserID    string            `json:"sub"`
    Username  string            `json:"username"`
    Email     string            `json:"email"`
    Roles     []string          `json:"roles"`
    Scopes    []string          `json:"scopes,omitempty"`
    IssuedAt  int64             `json:"iat"`
    ExpiresAt int64             `json:"exp"`
    Issuer    string            `json:"iss"`
    Audience  string            `json:"aud,omitempty"`
    JWTID     string            `json:"jti"`
    Custom    map[string]any `json:"custom,omitempty"`
}

type Session struct {
    ID           string    `json:"id"`
    UserID       string    `json:"user_id"`
    RefreshToken string    `json:"refresh_token"`
    UserAgent    string    `json:"user_agent"`
    IPAddress    string    `json:"ip_address"`
    LastActive   time.Time `json:"last_active"`
    ExpiresAt    time.Time `json:"expires_at"`
    CreatedAt    time.Time `json:"created_at"`
}

type Permission struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Resource    string    `json:"resource"` 
    Action      string    `json:"action"`   
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type RolePermission struct {
    RoleID       string    `json:"role_id"`
    PermissionID string    `json:"permission_id"`
    CreatedAt    time.Time `json:"created_at"`
}