package main

import (
	"flag"
	"fmt"
	"gomicro/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"golang.org/x/time/rate"

	"gomicro/util"
)

func main() {

	name := flag.String("name", "", "服务名称")
	port := flag.Int("p", 0, "服务端口")

	flag.Parse()

	if *name == "" {
		log.Fatal("请指定服务名称")
	}
	if *port == 0 {
		log.Fatal("请指定服务端口号")
	}

	util.SetServiceNameAndPort(*name, *port)

	user := services.UserService{}
	limiter := rate.NewLimiter(1, 5)
	endp := services.RateLimit(limiter)(services.GetUserEndpoint(&user))

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(services.MyErrorEncoder),
	}

	serverHandler := httptransport.NewServer(endp, services.DecodeUserRequest, services.EncodeUserResponse, options...)

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
		err := http.ListenAndServe(":"+strconv.Itoa(util.ServicePort), r)
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
