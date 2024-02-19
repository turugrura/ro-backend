package service

import "github.com/golang-jwt/jwt"

type AccessTokenResponse struct {
	AccessToken  string
	RefreshToken string
}

type AccessTokenRequest struct {
	UserAgent string
	UserId    string
}

type GenerateRefreshTokenRequest struct {
	UserAgent string
	UserId    string
}

type RefreshTokenRequest struct {
	Id        string
	UserAgent string
	UserId    string
	Count     uint32
}

type TokenService interface {
	DecodeToken(string) (*jwt.StandardClaims, error)
	GenerateAccessToken(AccessTokenRequest) (*AccessTokenResponse, error)
	GenerateRefreshToken(GenerateRefreshTokenRequest) (*string, error)
	RefreshToken(RefreshTokenRequest) (*AccessTokenResponse, error)
	RevokeTokenByUserId(string) error
}
