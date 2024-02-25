package repository

import "time"

type AuthenticationData struct {
	Id            string    `bson:"_id,omitempty"`
	AuthReference string    `bson:"auth_reference"`
	Email         string    `bson:"email"`
	CreatedAt     time.Time `bson:"created_at"`
}

type CreateAuthenticationDataInput struct {
	AuthReference string `json:"auth_reference"`
	Email         string `json:"email"`
}

type AuthenticationDataRepository interface {
	CreateAuthenticationData(CreateAuthenticationDataInput) (*AuthenticationData, error)
	GetAuthenticationById(string) (*AuthenticationData, error)
	GetAuthenticationByEmail(string) (*AuthenticationData, error)
	DeleteAuthenticationDataById(string) error
}
