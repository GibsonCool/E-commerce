package repositories

import (
	"E-commerce/datamodels"
	"github.com/jinzhu/gorm"
)

type IUser interface {
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
}

type UserRepository struct {
	sqlDb *gorm.DB
}

func (u *UserRepository) Select(userName string) (user *datamodels.User, err error) {
	user = &datamodels.User{}
	err = u.sqlDb.Where("user_name = ?", userName).First(user).Error
	return
}

func (u *UserRepository) Insert(user *datamodels.User) (userId int64, err error) {
	db := u.sqlDb.Create(user)
	userId = db.RowsAffected
	err = db.Error
	return
}

func NewUserRepository(sqlDb *gorm.DB) *UserRepository {
	return &UserRepository{sqlDb: sqlDb}
}
