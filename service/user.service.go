package service

import (
	"fmt"
	"ro-backend/repository"
)

type userService struct {
	userRepository repository.UserRepository
}

func (s userService) PatchUser(r PatchUserRequest) (*repository.User, error) {
	err := s.userRepository.PatchUser(r.Id, repository.UpdateUserInput{
		Name: r.Name,
	})
	if err != nil {
		return nil, err
	}

	return s.userRepository.FindUserById(r.Id)
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return userService{userRepository: userRepo}
}

func (s userService) CreateUser(req CreateUserRequest) (*repository.User, error) {
	existedUser, _ := s.userRepository.FindUserByEmail(req.Email)
	if existedUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	var user = repository.CreateUserInput{
		Name:            req.Name,
		Email:           req.Email,
		Role:            repository.UserRole.User,
		RegisterChannel: req.Channel,
	}
	return s.userRepository.CreateUser(user)
}

func (s userService) FindUserById(id string) (*repository.User, error) {
	return s.userRepository.FindUserById(id)
}

func (s userService) FindUserByEmail(email string) (*repository.User, error) {
	return s.userRepository.FindUserByEmail(email)
}
