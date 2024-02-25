package repository

import "time"

type AuthenticationData struct {
	Id        string    `bson:"_id,omitempty"`
	Code      string    `bson:"code"`
	Channel   string    `bson:"channel"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
}

type CreateAuthDataInput struct {
	Code      string    `bson:"code"`
	Channel   string    `bson:"channel"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
}

type DeleteAuthDataInput struct {
	Email string `json:"email"`
}

type PartialSearchAuthDataInput struct {
	Code  string `bson:"code,omitempty"`
	Email string `bson:"email,omitempty"`
}

type AuthenticationDataRepository interface {
	CreateAuthenticationData(CreateAuthDataInput) (*string, error)
	PartialSearchAuthData(PartialSearchAuthDataInput) (*AuthenticationData, error)
	DeleteAuthenticationDataById(string) error
	DeleteAuthenticationDataByEmail(string) error
}
