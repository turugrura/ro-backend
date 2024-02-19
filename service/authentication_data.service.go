package service

import "ro-backend/repository"

type authenticationDataService struct {
	repo repository.AuthenticationDataRepository
}

func NewAuthenticationDataService(repo repository.AuthenticationDataRepository) AuthenticationDataService {
	return authenticationDataService{repo: repo}
}

func (s authenticationDataService) FindAuthenticationDataByCode(code string) (*AuthenticationData, error) {
	authData, err := s.repo.GetAuthenticationById(code)
	if err != nil {
		return nil, err
	}

	return &AuthenticationData{
		AuthReference: authData.AuthReference,
		Code:          authData.Id,
		CreatedAt:     authData.CreatedAt,
		Email:         authData.Email,
	}, nil
}

func (s authenticationDataService) FindAuthenticationDataByEmail(email string) (*AuthenticationData, error) {
	authData, err := s.repo.GetAuthenticationByEmail(email)
	if err != nil {
		return nil, err
	}

	return &AuthenticationData{
		AuthReference: authData.AuthReference,
		Code:          authData.Id,
		CreatedAt:     authData.CreatedAt,
		Email:         authData.Email,
	}, nil
}

func (s authenticationDataService) CreateAuthenticationData(req AuthenticationDataRequest) (*AuthenticationData, error) {
	authData, err := s.repo.CreateAuthenticationData(repository.CreateAuthenticationDataInput{
		AuthReference: req.AuthReference,
		Email:         req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &AuthenticationData{
		AuthReference: authData.AuthReference,
		Code:          authData.Id,
		CreatedAt:     authData.CreatedAt,
		Email:         authData.Email,
	}, nil
}
