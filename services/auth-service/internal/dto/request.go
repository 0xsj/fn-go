package dto


type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	UserAgent string `json:"userAgent,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
}

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Phone     string `json:"phone,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// ValidateTokenRequest represents a token validation request
type ValidateTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// RevokeTokenRequest represents a token revocation request
type RevokeTokenRequest struct {
	Token string `json:"token" validate:"required"`
	TokenType string `json:"tokenType,omitempty"` // "access", "refresh", or "all"
}

// LogoutRequest represents a logout request
type LogoutRequest struct {
	UserID string `json:"userId" validate:"required"`
	SessionID string `json:"sessionId,omitempty"`
	RevokeAllSessions bool `json:"revokeAllSessions,omitempty"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	UserID string `json:"userId" validate:"required"`
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Token string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyEmailRequest represents an email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// CreatePermissionRequest represents a permission creation request
type CreatePermissionRequest struct {
	Name string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Resource string `json:"resource" validate:"required"`
	Action string `json:"action" validate:"required"`
}

// AssignPermissionRequest represents a role-permission assignment request
type AssignPermissionRequest struct {
	RoleID string `json:"roleId" validate:"required"`
	PermissionID string `json:"permissionId" validate:"required"`
}