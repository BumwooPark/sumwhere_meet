package profile

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"showper_server/factory"
	"time"
)

type Profile struct {
	ID        int64     `json:"id" xorm:"id pk autoincr"`
	UserID    int64     `json:"userID" xorm:"user_id" valid:"int,required"`
	Age       int       `json:"age" xorm:"age" valid:"int,range(10|60),required"`
	Job       string    `json:"job" xorm:"job" valid:"utfletter,required"`
	Phone     string    `json:"phone" xorm:"phone VARCHAR(30)" valid:"numeric"`
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
