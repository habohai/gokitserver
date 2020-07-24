package services

import (
	"context"
	"fmt"
	"gomicro/util"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"

	"github.com/hashicorp/vic/lib/apiservers/service/restapi/handlers/errors"
)

type UserRequest struct {
	UID    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

// RateLimit 加入限流功能的 中间件
func RateLimit(limiter *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if limiter.Allow() {
				return nil, errors.NewError(http.StatusTooManyRequests, "too many requests")
			}
			return next(ctx, request)
		}
	}
}

func GetUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
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
