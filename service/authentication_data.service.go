package service

import (
	"ro-backend/repository"
)

func NewAuthenticationDataService(repo repository.AuthenticationDataRepository) AuthenticationDataService {
	return authenticationDataService{repo: repo}
}

type authenticationDataService struct {
	repo repository.AuthenticationDataRepository
}

func (s authenticationDataService) DeleteAuthenticationDataByEmail(email string) error {
	return s.repo.DeleteAuthenticationDataByEmail(email)
}

func (s authenticationDataService) DeleteAuthenticationData(id string) error {
	return s.repo.DeleteAuthenticationDataById(id)
}

func (s authenticationDataService) FindAuthenticationDataByCode(code string) (*repository.AuthenticationData, error) {
	return s.repo.PartialSearchAuthData(repository.PartialSearchAuthDataInput{
		Code: code,
	})
}

func (s authenticationDataService) CreateAuthenticationData(req AuthenticationDataRequest) (*repository.AuthenticationData, error) {
	_, err := s.repo.CreateAuthenticationData(repository.CreateAuthDataInput{
		Channel: req.Channel,
		Email:   req.Email,
		Code:    req.Code,
	})
	if err != nil {
		return nil, err
	}

	return s.repo.PartialSearchAuthData(repository.PartialSearchAuthDataInput{Code: req.Code})
}
