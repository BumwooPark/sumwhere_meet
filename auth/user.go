package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"sumwhere_meet/factory"
	"time"
)

type User struct {
	ID        int64     `json:"id" xorm:"id pk autoincr" valid:"-"`
	Email     string    `json:"email" xorm:"email VARCHAR(30) unique not null" valid:"email"`
	Password  string    `json:"password" xorm:"VARCHAR(100) null password" valid:"stringlength(8|20),alphanum"`
	Nickname  string    `json:"nickname" xorm:"not null nickname" valid:"stringlength(2|10),required~유효하지않은 닉네임입니다."`
	Admin     bool      `json:"admin" xorm:"admin default false"`
	CreatedAt time.Time `json:"createAt" xorm:"created" valid:"-"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated" valid:"-"`
}

func (u *User) Insert(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Search(ctx context.Context) (*User, error) {
	var user User
	result, err := factory.DB(ctx).Where("email = ?", u.Email).Get(&user)
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, errors.New("")
	}

	return &user, nil
}

func (u *User) GetSocialUser(ctx context.Context) *User {
	var user User
	result, err := factory.DB(ctx).Where("email = ?", u.Email).Get(&user)
	if err != nil && !result {
		return nil
	}
	return &user
}

func (u *User) GenerateToken() (t string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = u.ID
	claims["email"] = u.Email
	claims["exp"] = time.Now().Add(time.Hour * 8760).Unix()

	t, err = token.SignedString([]byte("bumwoopark"))
	return
}
