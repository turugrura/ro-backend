package service

type UserResponse struct {
	Id        string
	Name      string
	Email     string
	Status    string
	Role      string
	CreatedAt string
}

type CreateUserRequest struct {
	Name  string
	Email string
}

type UserService interface {
	CreateUser(CreateUserRequest) (*UserResponse, error)
	FindUserById(string) (*UserResponse, error)
	FindUserByEmail(string) (*UserResponse, error)
}
