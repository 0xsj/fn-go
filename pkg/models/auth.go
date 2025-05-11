package models

import "time"

type AuthRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    Token        string    `json:"token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
    User         User      `json:"user"`
}

type TokenValidationRequest struct {
    Token string `json:"token" validate:"required"`
}

type TokenValidationResponse struct {
    Valid    bool   `json:"valid"`
    UserID   string `json:"user_id,omitempty"`
    Username string `json:"username,omitempty"`
    Role     Role   `json:"role,omitempty"`
}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}