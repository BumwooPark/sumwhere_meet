package auth

import (
	"context"
	"showper_server/factory"
	"time"
)

type SocialAuth struct {
	ID         int64     `json:"id" xorm:"id pk autoincr"`
	SocialID   string    `json:"socialId" xorm:"social_id"`
	UserID     int64     `json:"userId" xorm:"user_id unique"`
	SocialType string    `json:"socialType" xorm:"social_type VARCHAR(10)"`
	Token      string    `json:"token" xorm:"token VARCHAR(255)"`
	CreatedAt  time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt  time.Time `json:"updatedAt" xorm:"updated"`
}

func (s *SocialAuth) Create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(s)
	if err != nil {
		return err
	}
	return nil
}
