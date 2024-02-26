package handler

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

const (
	ErrUnverifiedEmail        = "Email is unverified"
	ErrEmptyEmail             = "Email is empty"
	ErrNotTimeForRefreshToken = "token is not valid yet"
	ErrUserNotFound           = "User not found"

	ErrForbidden        = "forbidden"
	ErrUnAuthentication = "unAuthentication"
	ErrNotMyPreset      = "not my preset"
)

func WriteErr(w http.ResponseWriter, msg string) {
	var message = msg
	var httpStatus = http.StatusInternalServerError

	switch msg {
	case ErrNotMyPreset:
		httpStatus = http.StatusNotFound
		message = http.StatusText(httpStatus)
	case mongo.ErrNoDocuments.Error():
		httpStatus = http.StatusNotFound
		message = http.StatusText(httpStatus)
	case primitive.ErrInvalidHex.Error():
		httpStatus = http.StatusBadRequest
	case ErrForbidden:
		httpStatus = http.StatusForbidden
		message = http.StatusText(httpStatus)
	case ErrUnAuthentication:
		httpStatus = http.StatusUnauthorized
		message = http.StatusText(httpStatus)
	}

	res := ErrorResponse{
		Message: message,
	}

	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(res)
}

func WriteOK(w http.ResponseWriter, res interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func WriteCreated(w http.ResponseWriter, res interface{}) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func WriteNoContent(w http.ResponseWriter, res interface{}) {
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(res)
}
