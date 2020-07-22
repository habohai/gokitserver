package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	mymux "github.com/gorilla/mux"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	vars := mymux.Vars(r)
	if uid, ok := vars["uid"]; ok {
		uid, _ := strconv.Atoi(uid)
		return UserRequest{
			UID:    uid,
			Method: r.Method,
		}, nil
	}

	return nil, errors.New("参数错误")
}

func EncodeUserResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
