package service

import "time"

type AuthenticationData struct {
	AuthReference string
	Code          string
	CreatedAt     time.Time
	Email         string
}

type AuthenticationDataRequest struct {
	AuthReference string
	Email         string
}

type AuthenticationDataService interface {
	CreateAuthenticationData(AuthenticationDataRequest) (*AuthenticationData, error)
	FindAuthenticationDataByCode(string) (*AuthenticationData, error)
	FindAuthenticationDataByEmail(string) (*AuthenticationData, error)
	DeleteAuthenticationData(string) error
}
