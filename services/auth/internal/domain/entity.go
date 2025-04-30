package domain

import "time"

type Token struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	TokenValue   string    `json:"tokenValue"` 
	TokenType    string    `json:"tokenType"`  
	ExpiresAt    time.Time `json:"expiresAt"`
	LastUsedAt   time.Time `json:"lastUsedAt,omitempty"`
	ClientIP     string    `json:"clientIp,omitempty"`
	UserAgent    string    `json:"userAgent,omitempty"`
	Revoked      bool      `json:"revoked"`
	RevokedAt    *time.Time `json:"revokedAt,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AccessRequest struct {
	ID          string    `json:"id"`
	CompanyName string    `json:"companyName"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	AccountType string    `json:"accountType"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PasswordReset struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId,omitempty"`
	Token     string    `json:"token,omitempty"`
	Used      bool      `json:"used"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type AccountAlert struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Action      string     `json:"action,omitempty"`
	Link        string     `json:"link,omitempty"`
	UserID      string     `json:"userId,omitempty"`
	EntityID    string     `json:"entityId,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type Session struct {
	ID           string     `json:"id"`
	UserID       string     `json:"userId"`
	RefreshToken string     `json:"refreshToken"` 
	ClientIP     string     `json:"clientIp,omitempty"`
	UserAgent    string     `json:"userAgent,omitempty"`
	LastActiveAt time.Time  `json:"lastActiveAt"`
	ExpiresAt    time.Time  `json:"expiresAt"`
	TerminatedAt *time.Time `json:"terminatedAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
}

type AuditLog struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId,omitempty"`
	Type      string    `json:"type,omitempty"` 
	IPAddress string    `json:"ipAddress,omitempty"`
	UserAgent string    `json:"userAgent,omitempty"`
	MessageID string    `json:"messageId,omitempty"`
	Keywords  string    `json:"keywords,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
