package services

type IUserService interface {
	GetName(userID int) string
}

type UserService struct{}

func (u *UserService) GetName(userID int) string {
	if userID == 101 {
		return "habo"
	}
	return "guest"
}
