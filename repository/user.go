package repository

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
	Id        string `bson:"_id,omitempty"`
	Name      string `bson:"name"`
	Email     string `bson:"email"`
	Status    string `bson:"status"`
	Role      string `bson:"role"`
	CreatedAt string `bson:"created_at"`
}

type CreateUserInput struct {
	Name  string
	Email string
	Role  string
}

type UserRepository interface {
	CreateUser(CreateUserInput) (*User, error)
	FindUserById(string) (*User, error)
	FindUserByEmail(string) (*User, error)
}
