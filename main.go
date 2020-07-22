package main

import (
	"fmt"
	"gomicro/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"

	"gomicro/util"
)

func main() {
	user := services.UserService{}
	endp := services.GetUserEndpoint(&user)

	serverHandler := httptransport.NewServer(endp, services.DecodeUserRequest, services.EncodeUserResponse)

	r := mymux.NewRouter()
	// r.Handle(`/user/{uid:\d+}`, serverHandler)
	{
		r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
		r.Methods("GET").Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			w.Write([]byte(`{"status": "ok"}`))
		})
	}

	errChan := make(chan error)

	go func() {
		util.RegisterService()
		err := http.ListenAndServe(":9050", r)
		if err != nil {
			log.Println(err)
			errChan <- err
		}
	}()

	go func() {
		sigC := make(chan os.Signal)
		signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigC)
	}()

	getErr := <-errChan
	util.DeregisterService()
	log.Println(getErr)
}
