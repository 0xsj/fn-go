// services/auth-service/pkg/jwt/jwt.go
package jwt

import (
	"time"

	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey            []byte
	accessTokenExpiry    time.Duration
	refreshTokenExpiry   time.Duration
	issuer               string
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, accessTokenExpiry, refreshTokenExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            []byte(secretKey),
		accessTokenExpiry:    accessTokenExpiry,
		refreshTokenExpiry:   refreshTokenExpiry,
		issuer:               "fn-go-auth-service",
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTManager) GenerateTokenPair(user *models.User) (accessToken, refreshToken string, err error) {
	// Generate access token
	accessToken, err = j.generateAccessToken(user)
	if err != nil {
		return "", "", domain.WithOperation(err, "generate_access_token")
	}

	// Generate refresh token
	refreshToken, err = j.generateRefreshToken(user)
	if err != nil {
		return "", "", domain.WithOperation(err, "generate_refresh_token")
	}

	return accessToken, refreshToken, nil
}

// generateAccessToken creates a short-lived access token
func (j *JWTManager) generateAccessToken(user *models.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"email":    user.Email,
		"roles":    []string{string(user.Role)},
		"iat":      now.Unix(),
		"exp":      now.Add(j.accessTokenExpiry).Unix(),
		"iss":      j.issuer,
		"jti":      uuid.New().String(),
		"type":     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", domain.NewInvalidTokenError()
	}

	return tokenString, nil
}

// generateRefreshToken creates a long-lived refresh token
func (j *JWTManager) generateRefreshToken(user *models.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"iat":  now.Unix(),
		"exp":  now.Add(j.refreshTokenExpiry).Unix(),
		"iss":  j.issuer,
		"jti":  uuid.New().String(),
		"type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", domain.NewInvalidTokenError()
	}

	return tokenString, nil
}

// ValidateAccessToken validates and parses an access token
func (j *JWTManager) ValidateAccessToken(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.NewInvalidTokenError()
		}
		return j.secretKey, nil
	})

	if err != nil {
		// Check if it's an expiration error by looking at the error message
		// This avoids importing the standard errors package
		if isTokenExpiredError(err) {
			return nil, domain.NewTokenExpiredError()
		}
		return nil, domain.NewInvalidTokenError()
	}

	if !token.Valid {
		return nil, domain.NewInvalidTokenError()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.NewInvalidTokenError()
	}

	// Verify token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, domain.NewInvalidTokenError()
	}

	// Extract and validate required claims
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, domain.NewInvalidTokenError()
	}

	username, _ := claims["username"].(string)
	email, _ := claims["email"].(string)
	jwtID, _ := claims["jti"].(string)

	// Extract roles
	var roles []string
	if rolesInterface, ok := claims["roles"].([]interface{}); ok {
		roles = make([]string, len(rolesInterface))
		for i, role := range rolesInterface {
			if roleStr, ok := role.(string); ok {
				roles[i] = roleStr
			}
		}
	}

	// Extract timestamps
	iat, _ := claims["iat"].(float64)
	exp, _ := claims["exp"].(float64)
	issuer, _ := claims["iss"].(string)

	tokenClaims := &models.TokenClaims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Roles:     roles,
		IssuedAt:  int64(iat),
		ExpiresAt: int64(exp),
		Issuer:    issuer,
		JWTID:     jwtID,
	}

	return tokenClaims, nil
}

// ValidateRefreshToken validates and parses a refresh token
func (j *JWTManager) ValidateRefreshToken(tokenString string) (userID string, jwtID string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.NewInvalidTokenError()
		}
		return j.secretKey, nil
	})

	if err != nil {
		if isTokenExpiredError(err) {
			return "", "", domain.NewTokenExpiredError()
		}
		return "", "", domain.NewInvalidTokenError()
	}

	if !token.Valid {
		return "", "", domain.NewInvalidTokenError()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", domain.NewInvalidTokenError()
	}

	// Verify token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", "", domain.NewInvalidTokenError()
	}

	// Extract required claims
	userID, ok = claims["sub"].(string)
	if !ok || userID == "" {
		return "", "", domain.NewInvalidTokenError()
	}

	jwtID, ok = claims["jti"].(string)
	if !ok || jwtID == "" {
		return "", "", domain.NewInvalidTokenError()
	}

	return userID, jwtID, nil
}

// ExtractTokenWithoutValidation extracts claims without validating the token
// Useful for debugging or when you need to inspect expired tokens
func (j *JWTManager) ExtractTokenWithoutValidation(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, domain.NewInvalidTokenError()
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, domain.NewInvalidTokenError()
}

func (j *JWTManager) GetTokenExpiry() (accessExpiry, refreshExpiry time.Duration) {
	return j.accessTokenExpiry, j.refreshTokenExpiry
}

func isTokenExpiredError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return errMsg == "token is expired" || errMsg == "Token is expired"
}