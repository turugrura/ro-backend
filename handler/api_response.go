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
	ErrForbidden              = "forbidden"
	ErrUnAuthentication       = "unAuthentication"
	ErrUnverifiedEmail        = "Email is unverified"
	ErrEmptyEmail             = "Email is empty"
	ErrUserNotFound           = "User not found"
	ErrNotTimeForRefreshToken = "token is not valid yet"
)

func WriteErr(w http.ResponseWriter, msg string) {
	var message = msg
	var httpStatus = http.StatusInternalServerError

	switch msg {
	case mongo.ErrNoDocuments.Error():
		httpStatus = http.StatusNotFound
		message = http.StatusText(httpStatus)
	case primitive.ErrInvalidHex.Error():
		httpStatus = http.StatusBadRequest
		message = http.StatusText(httpStatus)
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
