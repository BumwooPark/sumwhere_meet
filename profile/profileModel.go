package profile

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"sumwhere_meet/factory"
	"time"
)

type Profile struct {
	ID        int64     `json:"id" xorm:"id pk autoincr" example:"1"`
	UserID    int64     `json:"userID" xorm:"user_id" valid:"int,required" example:"1"`
	Age       int       `json:"age" xorm:"age" valid:"int,range(10|60),required" example:"30"`
	Job       string    `json:"job" xorm:"job" valid:"utfletter,required" example:"공무원"`
	Phone     string    `json:"phone" xorm:"phone VARCHAR(30)" valid:"numeric" example:"01051416906"`
	Gender    string    `json:"gender" xorm:"gender VARCHAR(10)" valid:"required" example:"남성"`
	Area      int       `json:"area" xorm:"area" valid:"requried"`
	Interest  int       `json:"interest" xorm:"interest" valid:"requried"`
	CreatedAt time.Time `json:"createAt" xorm:"created"`
	UpdatedAt time.Time `json:"updateAt" xorm:"updated"`
}

func (p *Profile) Create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Profile) Get(ctx context.Context, pk int64) error {

	result, err := factory.DB(ctx).ID(pk).Get(p)
	if err != nil {
		return err
	}

	if !result {
		return errors.New(fmt.Sprintf("user %d does not exist", pk))
	}
	return nil
}

func (p *Profile) Update(ctx context.Context) error {
	_, err := factory.DB(ctx).ID(p.ID).Update(p)
	if err != nil {
		return err
	}
	return nil
}
