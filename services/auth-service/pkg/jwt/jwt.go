package jwt

import (
	"time"

	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secretKey            []byte
	accessTokenExpiry    time.Duration
	refreshTokenExpiry   time.Duration
	issuer               string
}

func NewJWTManager(secretKey string, accessTokenExpiry, refreshTokenExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            []byte(secretKey),
		accessTokenExpiry:    accessTokenExpiry,
		refreshTokenExpiry:   refreshTokenExpiry,
		issuer:               "fn-go-auth-service",
	}
}

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