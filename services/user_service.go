package services

import (
	"E-commerce/datamodels"
	"E-commerce/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepository repositories.IUser
}

func NewUserService(userRepository repositories.IUser) IUserService {
	return &UserService{UserRepository: userRepository}
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPwd)
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	password, errPwd := GeneratePassword(user.HashPwd)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPwd = string(password)
	return u.UserRepository.Insert(user)
}

// bcrypt 使用随机盐，生成对应密码 hash 值
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// 使用密码以及上面产生的密码hash值进行校验
// 成功返回 nil  失败返回错误 err
func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil

}
