package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
)

type IRequest interface {
	Validate() error
}

func ParseBody[T IRequest](r *http.Request) (*T, error) {
	var req T
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	if err = req.Validate(); err != nil {
		return &req, err
	}

	return &req, nil
}

func ParseParam[T IRequest](r *http.Request) (*T, error) {
	var req T
	schema.NewDecoder().Decode(&req, r.URL.Query())
	if err := req.Validate(); err != nil {
		return &req, err
	}

	return &req, nil
}

func Parse[T IRequest](r *http.Request) (*T, error) {
	var req T
	schema.NewDecoder().Decode(&req, r.URL.Query())
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	if err = req.Validate(); err != nil {
		return &req, err
	}

	return &req, nil
}

type Response[T any] struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func NewResponse[T any](code int, message string, err string, data T) *Response[T] {
	return &Response[T]{
		Code:    code,
		Message: message,
		Error:   err,
		Data:    data,
	}
}

func (r *Response[T]) ToJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	json.NewEncoder(w).Encode(r)
}
