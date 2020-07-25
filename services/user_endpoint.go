package services

import (
	"context"
	"fmt"
	"gomicro/util"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
)

type UserRequest struct {
	UID    int `json:"uid"`
	Method string
	Token  string
}

type UserResponse struct {
	Result string `json:"result"`
}

// CheckTokenMiddleware token验证中间件
func CheckTokenMiddleware() endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest) //通过类型断言获取请求结构体
			uc := UserClaim{}
			//下面的r.Token是在代码DecodeUserRequest那里封装进去的
			getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (i interface{}, e error) {
				return []byte(secKey), err
			})
			fmt.Println(err, 123)
			if getToken != nil && getToken.Valid { //验证通过
				newCtx := context.WithValue(ctx, "LoginUser", getToken.Claims.(*UserClaim).Uname)
				return next(newCtx, request)
			} else {
				return nil, util.NewMyError(403, "error token")
			}

			//logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)

		}
	}
}

// UserServiceLogMiddleware 日志 中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			logger.Log("method", r.Method, "event", "get user", "user ID", r.UID)
			return next(ctx, request)
		}
	}
}

// RateLimit 加入限流功能的 中间件
func RateLimit(limiter *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limiter.Allow() {
				return nil, util.NewMyError(http.StatusTooManyRequests, "too many requests")
			}
			return next(ctx, request)
		}
	}
}

func GetUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		fmt.Println("当前登录用户名：", ctx.Value("LoginUser"))
		result := "nothing"
		if r.Method == "GET" {
			result = userService.GetName(r.UID) + strconv.Itoa(util.ServicePort)
		} else if r.Method == "DELETE" {
			if err := userService.DelUser(r.UID); err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("uid为%d的用户删除成功", r.UID)
			}
		}

		return UserResponse{
			Result: result,
		}, nil
	}
}
