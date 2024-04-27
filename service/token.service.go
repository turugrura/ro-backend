package service

import (
	"fmt"
	"ro-backend/configuration"
	"ro-backend/repository"
	"time"

	"github.com/golang-jwt/jwt"
)

type tokenService struct {
	repo repository.RefreshTokenRepository
}

func NewTokenService(repo repository.RefreshTokenRepository) TokenService {
	return tokenService{repo: repo}
}

func (s tokenService) GenerateAccessToken(req AccessTokenRequest) (*AccessTokenResponse, error) {
	signedAccessToken, err := s.signAccessToken(req)
	if err != nil {
		return nil, err
	}

	signedRefreshToken, err := s.GenerateRefreshToken(GenerateRefreshTokenRequest{AccessTokenRequest: req})
	if err != nil {
		return nil, err
	}

	accessToken := AccessTokenResponse{AccessToken: *signedAccessToken, RefreshToken: *signedRefreshToken}

	return &accessToken, nil
}

func (s tokenService) GenerateRefreshToken(req GenerateRefreshTokenRequest) (*string, error) {
	var newRefreshToken = repository.CreateRefreshTokenInput{
		UserId:    req.UserId,
		Count:     1,
		UserAgent: req.UserAgent,
	}

	createdRefreshToken, err := s.repo.CreateRefreshToken(newRefreshToken)
	if err != nil {
		return nil, err
	}

	signedRefreshToken, err := s.signRefreshToken(createdRefreshToken.Id, req.UserId, newRefreshToken.Count)
	if err != nil {
		return nil, err
	}

	return signedRefreshToken, nil
}

func (s tokenService) RefreshToken(req RefreshTokenRequest) (*AccessTokenResponse, error) {
	token, err := s.repo.GetRefreshTokenById(req.Id)
	if err != nil {
		return nil, err
	}

	if token.Count != req.Count || token.UserAgent != req.UserAgent {
		s.RevokeTokenByUserId(req.UserId)
		return nil, fmt.Errorf("token is mismatch")
	}

	signedAccessToken, err := s.signAccessToken(req.AccessTokenRequest)
	if err != nil {
		return nil, err
	}

	nextCount := req.Count + 1
	signedRefreshToken, err := s.signRefreshToken(token.Id, req.UserId, nextCount)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.UpdateRefreshToken(repository.UpdateRefreshTokenInput{
		Id:    req.Id,
		Count: nextCount,
	})
	if err != nil {
		return nil, err
	}

	accessToken := AccessTokenResponse{AccessToken: *signedAccessToken, RefreshToken: *signedRefreshToken}

	return &accessToken, nil
}

func (s tokenService) RevokeTokenByUserId(userId string) error {
	return s.repo.DeleteRefreshTokenByUserId(userId)
}

type AppClaims struct {
	jwt.StandardClaims
	Role string
}

func (s tokenService) signAccessToken(r AccessTokenRequest) (*string, error) {
	var tokenPeriod = time.Duration(configuration.Config.Jwt.AccessTokenPeriodInMinutes)

	now := time.Now()
	claims := jwt.StandardClaims{
		Id:        r.UserId,
		Subject:   r.Role,
		Issuer:    r.Name,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Minute * tokenPeriod).Unix(),
	}

	return s.signToken(claims)
}

func (s tokenService) signRefreshToken(tokenId, userId string, count uint32) (*string, error) {
	var tokenPeriodInDays = time.Duration(configuration.Config.Jwt.RefreshTokenPeriodInDays)
	var rfNBF = time.Duration(configuration.Config.Jwt.RefreshTokenNotBeforeInMinutes)

	now := time.Now()
	claims := jwt.StandardClaims{
		Id:        tokenId,
		Subject:   fmt.Sprint(count),
		Issuer:    userId,
		IssuedAt:  now.Unix(),
		NotBefore: now.Add(time.Minute * rfNBF).Unix(),
		ExpiresAt: now.Add(time.Hour * 24 * tokenPeriodInDays).Unix(),
	}

	return s.signToken(claims)
}

func (s tokenService) signToken(claims jwt.Claims) (*string, error) {
	var jwtSecret = configuration.Config.Jwt.Secret

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}

func (s tokenService) DecodeToken(tokenString string) (*jwt.StandardClaims, error) {
	var jwtSecret = configuration.Config.Jwt.Secret

	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return &claims, nil
}
