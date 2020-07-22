package main

import (
	"gomicro/services"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
)

func main() {
	user := services.UserService{}
	endp := services.GetUserEndpoint(&user)

	serverHandler := httptransport.NewServer(endp, services.DecodeUserRequest, services.EncodeUserResponse)

	r := mymux.NewRouter()
	// r.Handle(`/user/{uid:\d+}`, serverHandler)
	r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)

	http.ListenAndServe(":8080", r)
}
