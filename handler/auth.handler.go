package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ro-backend/appError"
	"ro-backend/configuration"
	"ro-backend/core"
	"ro-backend/repository"
	"ro-backend/service"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	Logout(http.ResponseWriter, *http.Request)
	RefreshToken(http.ResponseWriter, *http.Request)
	AuthenticationCallback(http.ResponseWriter, *http.Request)
}

type AuthHandlerParam struct {
	UserService               service.UserService
	AuthenticationDataService service.AuthenticationDataService
	TokenService              service.TokenService
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

func NewAuthHandler(param AuthHandlerParam) AuthHandler {
	return authHandler{userService: param.UserService, authenticationDataService: param.AuthenticationDataService, tokenService: param.TokenService}
}

func (h authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var p LoginRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	if p.AuthorizationCode == "" {
		return
	}

	authData, err := h.authenticationDataService.FindAuthenticationDataByCode(p.AuthorizationCode)
	if err == mongo.ErrNoDocuments {
		core.WriteErr(w, appError.ErrUnAuthentication)
		return
	}
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	user, err := h.userService.FindUserByEmail(authData.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		core.WriteErr(w, err.Error())
		return
	}

	if user == nil {
		core.WriteErr(w, appError.ErrUserNotFound)
		return
	}

	generatedToken, err := h.tokenService.GenerateAccessToken(service.AccessTokenRequest{
		UserId:    user.Id,
		UserAgent: r.UserAgent(),
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		Role:      user.Role,
	})
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	err = h.authenticationDataService.DeleteAuthenticationDataByEmail(authData.Email)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	var loginResponse = LoginResponse{
		AccessToken:  generatedToken.AccessToken,
		RefreshToken: generatedToken.RefreshToken,
	}
	core.WriteOK(w, loginResponse)
}

func (h authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	err := h.tokenService.RevokeTokenByUserId(userId)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, nil)
}

func (h authHandler) AuthenticationCallback(w http.ResponseWriter, r *http.Request) {
	userInfo, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		core.WriteErr(w, err.Error())
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
		core.WriteErr(w, appError.ErrUnverifiedEmail)
		return
	}

	if userInfo.Email == "" {
		core.WriteErr(w, appError.ErrEmptyEmail)
		return
	}

	provider := mux.Vars(r)["provider"]

	user, err := h.userService.FindUserByEmail(userInfo.Email)
	if err == mongo.ErrNoDocuments {
		name := strings.Split(userInfo.Email, "@")[0]

		user, err = h.userService.CreateUser(service.CreateUserRequest{
			Name:    name,
			Email:   userInfo.Email,
			Channel: provider,
		})

		if err != nil {
			core.WriteErr(w, err.Error())
			return
		}
	} else if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	code := uuid.NewString()
	_, err = h.authenticationDataService.CreateAuthenticationData(service.AuthenticationDataRequest{
		Channel: provider,
		Email:   user.Email,
		Code:    code,
	})
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	var redirectUrl = fmt.Sprintf("%v?%v=%v", configuration.Config.Auth.PostAuthenticationRedirectUrl, "auth_code", code)

	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

func (h authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var p RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	claims, err := h.tokenService.DecodeToken(p.RefreshToken)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	count, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	user, err := h.userService.FindUserById(claims.Issuer)
	if err == mongo.ErrNoDocuments {
		core.WriteErr(w, appError.ErrUnAuthentication)
		return
	}
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}
	if user.Status != repository.UserStatus.Active {
		core.WriteErr(w, appError.ErrUserInactive)
		return
	}

	newToken, err := h.tokenService.RefreshToken(service.RefreshTokenRequest{
		Id:    claims.Id,
		Count: uint32(count),
		AccessTokenRequest: service.AccessTokenRequest{
			UserAgent: r.UserAgent(),
			UserId:    user.Id,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			Role:      user.Role,
		},
	})
	if err == mongo.ErrNoDocuments {
		core.WriteErr(w, appError.ErrUnAuthentication)
		return
	}
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	var loginResponse = LoginResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	}

	core.WriteOK(w, loginResponse)
}
