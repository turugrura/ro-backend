package handler

import (
	"net/http"
	"ro-backend/service"
)

type userHandler struct {
	userService service.UserService
}

type UserHandler interface {
	GetMyProfile(http.ResponseWriter, *http.Request)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return userHandler{userService: userService}
}

type GetMyProfileResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

func (h userHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	user, err := h.userService.FindUserById(userId)
	if err != nil {
		WriteErr(w, ErrUnAuthentication)
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

	WriteOK(w, response)
}
