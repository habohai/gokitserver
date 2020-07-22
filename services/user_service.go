package services

import "errors"

type IUserService interface {
	GetName(userID int) string
	DelUser(userID int) error
}

type UserService struct{}

func (u *UserService) GetName(userID int) string {
	if userID == 101 {
		return "habo"
	}
	return "guest"
}

func (u *UserService) DelUser(userID int) error {
	if userID == 101 {
		return errors.New("无权限")
	}
	return nil
}
