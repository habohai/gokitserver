package services

import (
	"context"
	"fmt"
	"gomicro/util"
	"strconv"

	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct {
	UID    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
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
