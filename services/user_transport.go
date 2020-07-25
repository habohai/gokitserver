package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/haibeihabo/gokitserver/util"

	mymux "github.com/gorilla/mux"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	vars := mymux.Vars(r)
	if uid, ok := vars["uid"]; ok {
		uid, _ := strconv.Atoi(uid)
		return UserRequest{
			UID:    uid,
			Method: r.Method,
			Token:  r.URL.Query().Get("token"),
		}, nil
	}

	return nil, errors.New("参数错误")
}

func EncodeUserResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MyErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("content-type", contentType)

	if myerr, ok := err.(*util.MyError); ok {
		w.WriteHeader(myerr.Code)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(body)
}
