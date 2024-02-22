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
		fmt.Println(err)
		return
	}

	authorizationCode := p.AuthorizationCode
	fmt.Println("authorizationCode", authorizationCode)

	if authorizationCode == "" {
		return
	}

	authData, err := h.authenticationDataService.FindAuthenticationDataByCode(authorizationCode)
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := h.userService.FindUserByEmail(authData.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println(err)
		return
	}

	if user == nil {
		fmt.Println("User not found")
		return
	}

	generatedToken, err := h.tokenService.GenerateAccessToken(service.AccessTokenRequest{
		UserId:    user.Id,
		UserAgent: r.Header.Get("user-agent"),
	})
	if err != nil {
		fmt.Print(err)
		return
	}

	var loginResponse = LoginResponse{
		AccessToken:  generatedToken.AccessToken,
		RefreshToken: generatedToken.RefreshToken,
	}

	json.NewEncoder(w).Encode(loginResponse)
}

func (h authHandler) AuthenticationCallback(w http.ResponseWriter, r *http.Request) {
	userInfo, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
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
		fmt.Print("Email is unverified")
		return
	}

	if userInfo.Email == "" {
		fmt.Println("Email is empty")
		return
	}

	user, err := h.userService.FindUserByEmail(userInfo.Email)
	if err == mongo.ErrNoDocuments {
		user, err = h.userService.CreateUser(service.CreateUserRequest{
			Name:  "Unknown",
			Email: userInfo.Email,
		})

		if err != nil {
			fmt.Println(err)
			return
		}
	} else if err != nil {
		fmt.Println(err)
		return
	}

	createdAuthData, err := h.authenticationDataService.CreateAuthenticationData(service.AuthenticationDataRequest{
		AuthReference: "google",
		Email:         user.Email,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	var redirectUrl = fmt.Sprintf("%v?%v=%v", configuration.Config.Auth.PostAuthenticationRedirectUrl, "auth_code", createdAuthData.Code)

	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

func (h authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var p RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read body, %v\n", err), http.StatusBadRequest)
		return
	}

	claims, err := h.tokenService.DecodeToken(p.RefreshToken)
	if err != nil {
		fmt.Printf("cannot decode token, %v\n", err)
		return
	}

	count, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		fmt.Printf("cannot ParseUint, %v\n", err)
		return
	}

	newToken, err := h.tokenService.RefreshToken(service.RefreshTokenRequest{
		Id:        claims.Id,
		UserAgent: r.UserAgent(),
		UserId:    claims.Issuer,
		Count:     uint32(count),
	})
	if err != nil {
		fmt.Printf("cannot refresh token, %v\n", err)
		return
	}

	var loginResponse = LoginResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	}

	json.NewEncoder(w).Encode(loginResponse)
}
