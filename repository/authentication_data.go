package repository

type AuthenticationData struct {
	Id            string `bson:"_id,omitempty"`
	AuthReference string `bson:"auth_reference"`
	CreatedAt     string `bson:"created_at"`
	Email         string `bson:"email"`
}

type CreateAuthenticationDataInput struct {
	AuthReference string `json:"auth_reference"`
	Email         string `json:"email"`
}

type AuthenticationDataRepository interface {
	CreateAuthenticationData(CreateAuthenticationDataInput) (*AuthenticationData, error)
	GetAuthenticationById(string) (*AuthenticationData, error)
	GetAuthenticationByEmail(string) (*AuthenticationData, error)
}
