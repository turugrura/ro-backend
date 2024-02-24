package service

type AuthenticationData struct {
	AuthReference string
	Code          string
	CreatedAt     string
	Email         string
}

type AuthenticationDataRequest struct {
	AuthReference string
	Email         string
}

type AuthenticationDataService interface {
	CreateAuthenticationData(AuthenticationDataRequest) (*AuthenticationData, error)
	FindAuthenticationDataByCode(string) (*AuthenticationData, error)
	FindAuthenticationDataByEmail(string) (*AuthenticationData, error)
	DeleteAuthenticationData(string) error
}
