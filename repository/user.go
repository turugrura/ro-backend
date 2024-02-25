package repository

import "time"

type Status struct {
	Active   string
	InActive string
}

var UserStatus = Status{
	Active:   "active",
	InActive: "inactive",
}

type Role struct {
	Admin string
	User  string
}

var UserRole = Role{
	Admin: "admin",
	User:  "user",
}

type User struct {
	Id        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Status    string    `bson:"status"`
	Role      string    `bson:"role"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type CreateUserInput struct {
	Name  string
	Email string
	Role  string
}

type UpdateUserInput struct {
	Name      string    `bson:"name,omitempty"`
	Status    string    `bson:"status,omitempty"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type UserRepository interface {
	CreateUser(CreateUserInput) (*User, error)
	PatchUser(id string, u UpdateUserInput) error
	FindUserById(string) (*User, error)
	FindUserByEmail(string) (*User, error)
}
