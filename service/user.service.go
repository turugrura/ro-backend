package service

import (
	"fmt"
	"ro-backend/repository"
	"strings"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) userService {
	return userService{userRepository: userRepo}
}

func (s userService) CreateUser(req CreateUserRequest) (*UserResponse, error) {
	existedUser, _ := s.userRepository.FindUserByEmail(req.Email)
	if existedUser != nil {
		return nil, fmt.Errorf("email '%v' already registered", req.Email)
	}

	name := strings.Split(req.Email, "@")[0]

	var user = repository.CreateUserInput{
		Name:  name,
		Email: req.Email,
		Role:  repository.UserRole.User,
	}
	newUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		Id:        newUser.Id,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Status:    newUser.Status,
		Role:      newUser.Role,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

func (s userService) FindUserById(id string) (*UserResponse, error) {
	user, err := s.userRepository.FindUserById(id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s userService) FindUserByEmail(email string) (*UserResponse, error) {
	user, err := s.userRepository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}
