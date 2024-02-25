package service

import (
	"ro-backend/repository"
)

type AuthenticationDataRequest struct {
	Channel string
	Email   string
	Code    string
}

type AuthenticationDataService interface {
	CreateAuthenticationData(AuthenticationDataRequest) (*repository.AuthenticationData, error)
	FindAuthenticationDataByCode(string) (*repository.AuthenticationData, error)
	DeleteAuthenticationData(string) error
	DeleteAuthenticationDataByEmail(string) error
}
