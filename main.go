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

	kitlog "github.com/go-kit/kit/log"

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

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stdout)
	logger = kitlog.WithPrefix(logger, "mykit", "1.0")
	logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)

	user := services.UserService{} // 用户服务
	limiter := rate.NewLimiter(1, 5)
	endp := services.RateLimit(limiter)(services.UserServiceLogMiddleware(logger)(services.CheckTokenMiddleware()(services.GetUserEndpoint(&user))))

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(services.MyErrorEncoder),
	}

	serverHandler := httptransport.NewServer(endp, services.DecodeUserRequest, services.EncodeUserResponse, options...)

	// 增加 handle 用于获取用户token
	accessService := services.AccessService{}
	accessServiceEndpoint := services.AccessEndpoint(&accessService)
	accessHandler := httptransport.NewServer(accessServiceEndpoint, services.DecodeAccessRequest, services.EncodeAccessResponse, options...)

	r := mymux.NewRouter()
	// r.Handle(`/user/{uid:\d+}`, serverHandler)
	{
		r.Methods("POST").Path("/access-token").Handler(accessHandler) //注册token获取的handler
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
