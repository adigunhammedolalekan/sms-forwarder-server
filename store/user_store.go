package store

import (
	"errors"
	"github.com/adigunhammedolalekan/sms-forwarder/fn"
	"github.com/adigunhammedolalekan/sms-forwarder/types"
	"github.com/jinzhu/gorm"
)
//go:generate mockgen -destination=../mocks/user_store_mock.go -package=mocks github.com/adigunhammedolalekan/sms-forwarder/store UserStore
type UserStore interface {
	CreateUser(email, password string) (*types.User, error)
	AuthenticateUser(email, password string) (*types.User, error)
	FindUser(email string) (*types.User, error)
	FindUserById(id uint) (*types.User, error)
}

type databaseUserStore struct {
	db *gorm.DB
	tokenGenerator *fn.JwtTokenGenerator
}

func (d *databaseUserStore) CreateUser(email, password string) (*types.User, error) {
	if u, err := d.FindUser(email); err == nil && u.ID > 0 {
		return nil, errors.New("a user is already using that email")
	}
	user := types.NewUser(email, password)
	user.Password = fn.HashPassword(user.Password)
	if err := d.db.Create(user).Error; err != nil {
		return nil, err
	}
	user.Token = d.tokenGenerator.SignJwtToken(user.ID, user.Email)
	return user, nil
}

func (d *databaseUserStore) AuthenticateUser(email, password string) (*types.User, error) {
	user, err := d.FindUser(email)
	if err != nil || user.ID <= 0 {
		return nil, errors.New("invalid login credentials")
	}
	if ok := fn.VerifyPassword(user.Password, password); !ok {
		return nil, errors.New("invalid login credentials")
	}
	user.Token = d.tokenGenerator.SignJwtToken(user.ID, user.Email)
	return user, nil
}

func (d *databaseUserStore) FindUser(email string) (*types.User, error) {
	user := &types.User{}
	err := d.db.Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *databaseUserStore) FindUserById(id uint) (*types.User, error) {
	user := &types.User{}
	err := d.db.Table("users").Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserStore(db *gorm.DB, tokenGen *fn.JwtTokenGenerator) UserStore {
	return &databaseUserStore{db:db, tokenGenerator: tokenGen}
}

