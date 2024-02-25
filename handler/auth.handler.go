package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ro-backend/configuration"
	"ro-backend/service"
	"strconv"

	"github.com/markbates/goth/gothic"
	"go.mongodb.org/mongo-driver/mongo"
)

type authHandler struct {
	userService               service.UserService
	authenticationDataService service.AuthenticationDataService
	tokenService              service.TokenService
}

type AuthHandler interface {
	Login(http.ResponseWriter, *http.Request)
	RefreshToken(http.ResponseWriter, *http.Request)
	AuthenticationCallback(http.ResponseWriter, *http.Request)
}

type AuthHandlerParam struct {
	UserService               service.UserService
	AuthenticationDataService service.AuthenticationDataService
	TokenService              service.TokenService
}

func NewAuthHandler(param AuthHandlerParam) AuthHandler {
	return authHandler{userService: param.UserService, authenticationDataService: param.AuthenticationDataService, tokenService: param.TokenService}
}

type LoginRequest struct {
	AuthorizationCode string `json:"authorizationCode"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var p LoginRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	if p.AuthorizationCode == "" {
		return
	}

	authData, err := h.authenticationDataService.FindAuthenticationDataByCode(p.AuthorizationCode)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	user, err := h.userService.FindUserByEmail(authData.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		WriteErr(w, err.Error())
		return
	}

	if user == nil {
		WriteErr(w, ErrUserNotFound)
		return
	}

	generatedToken, err := h.tokenService.GenerateAccessToken(service.AccessTokenRequest{
		UserId:    user.Id,
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	err = h.authenticationDataService.DeleteAuthenticationData(p.AuthorizationCode)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	var loginResponse = LoginResponse{
		AccessToken:  generatedToken.AccessToken,
		RefreshToken: generatedToken.RefreshToken,
	}
	WriteOK(w, loginResponse)
}

func (h authHandler) AuthenticationCallback(w http.ResponseWriter, r *http.Request) {
	userInfo, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	// fmt.Println("UserID", userInfo.UserID)
	// fmt.Println("Email", userInfo.Email)
	// fmt.Println("AccessToken", userInfo.AccessToken)
	// fmt.Println("RefreshToken", userInfo.RefreshToken)
	// fmt.Println("ExpiresAt", userInfo.ExpiresAt)
	// fmt.Println("AvatarURL", userInfo.AvatarURL)
	// fmt.Println("verified_email", userInfo.RawData["verified_email"])
	// fmt.Println(user)

	if userInfo.RawData["verified_email"] == false {
		WriteErr(w, ErrUnverifiedEmail)
		return
	}

	if userInfo.Email == "" {
		WriteErr(w, ErrEmptyEmail)
		return
	}

	user, err := h.userService.FindUserByEmail(userInfo.Email)
	if err == mongo.ErrNoDocuments {
		user, err = h.userService.CreateUser(service.CreateUserRequest{
			Name:  "Unknown",
			Email: userInfo.Email,
		})

		if err != nil {
			WriteErr(w, err.Error())
			return
		}
	} else if err != nil {
		WriteErr(w, err.Error())
		return
	}

	createdAuthData, err := h.authenticationDataService.CreateAuthenticationData(service.AuthenticationDataRequest{
		AuthReference: "google",
		Email:         user.Email,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	var redirectUrl = fmt.Sprintf("%v?%v=%v", configuration.Config.Auth.PostAuthenticationRedirectUrl, "auth_code", createdAuthData.Code)

	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

func (h authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var p RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	claims, err := h.tokenService.DecodeToken(p.RefreshToken)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	count, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	newToken, err := h.tokenService.RefreshToken(service.RefreshTokenRequest{
		Id:        claims.Id,
		UserAgent: r.UserAgent(),
		UserId:    claims.Issuer,
		Count:     uint32(count),
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	var loginResponse = LoginResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	}

	WriteOK(w, loginResponse)
}
