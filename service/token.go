package service

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type AccessTokenResponse struct {
	AccessToken  string
	RefreshToken string
}

type AccessTokenRequest struct {
	UserAgent string
	UserId    string
	Name      string
	CreatedAt time.Time
	Role      string
}

type GenerateRefreshTokenRequest struct {
	AccessTokenRequest
}

type RefreshTokenRequest struct {
	AccessTokenRequest
	Id    string
	Count uint32
}

type TokenService interface {
	DecodeToken(string) (*jwt.StandardClaims, error)
	GenerateAccessToken(AccessTokenRequest) (*AccessTokenResponse, error)
	GenerateRefreshToken(GenerateRefreshTokenRequest) (*string, error)
	RefreshToken(RefreshTokenRequest) (*AccessTokenResponse, error)
	RevokeTokenByUserId(string) error
}
