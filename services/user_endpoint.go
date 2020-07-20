package services

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct {
	UID int `json:"uid"`
}

type UserResponse struct {
	Result string `json:"result"`
}

func GetUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := userService.GetName(r.UID)
		return UserResponse{
			Result: result,
		}, nil
	}
}
