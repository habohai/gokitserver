package main

import (
	"gomicro/services"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	user := services.UserService{}
	endp := services.GetUserEndpoint(&user)

	serverHandler := httptransport.NewServer(endp, services.DecodeUserRequest, services.EncodeUserResponse)
	http.ListenAndServe(":8080", serverHandler)
}
