package handler

import (
	"encoding/json"
	"net/http"
	"ro-backend/appError"
	"ro-backend/service"
	"time"
)

type UserHandler interface {
	GetMyProfile(http.ResponseWriter, *http.Request)
	PatchMyProfile(http.ResponseWriter, *http.Request)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return userHandler{userService: userService}
}

type userHandler struct {
	userService service.UserService
}

type GetMyProfileResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PatchMyProfileRequest struct {
	Name string `json:"name"`
}

func (h userHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	user, err := h.userService.FindUserById(userId)
	if err != nil {
		WriteErr(w, appError.ErrUnAuthentication)
		return
	}

	var response = GetMyProfileResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	WriteOK(w, response)
}

func (h userHandler) PatchMyProfile(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	var d PatchMyProfileRequest
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	user, err := h.userService.PatchUser(service.PatchUserRequest{
		Id:   userId,
		Name: d.Name,
	})
	if err != nil {
		WriteErr(w, err.Error())
		return
	}

	var response = GetMyProfileResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	WriteOK(w, response)
}
