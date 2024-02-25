package service

import (
	"ro-backend/repository"
)

type CreateUserRequest struct {
	Name  string
	Email string
}

type PatchUserRequest struct {
	Id   string
	Name string
}

type UserService interface {
	CreateUser(CreateUserRequest) (*repository.User, error)
	PatchUser(PatchUserRequest) (*repository.User, error)
	FindUserById(string) (*repository.User, error)
	FindUserByEmail(string) (*repository.User, error)
}
