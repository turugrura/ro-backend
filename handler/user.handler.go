package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ro-backend/service"
)

type userHandler struct {
	userService service.UserService
}

type UserHandler interface {
	GetMyProfile(res http.ResponseWriter, req *http.Request)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return userHandler{userService: userService}
}

type GetMyProfileResponse struct {
	Id        string
	Name      string
	Email     string
	Status    string
	Role      string
	CreatedAt string
}

func (h userHandler) GetMyProfile(res http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("userId")

	user, err := h.userService.FindUserById(userId)
	if err != nil {
		fmt.Println("cannot find user")
		return
	}

	var response = GetMyProfileResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	json.NewEncoder(res).Encode(response)
}
