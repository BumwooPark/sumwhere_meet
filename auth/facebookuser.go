package auth

import (
	"context"
	"fmt"
	"sumwhere_meet/factory"
)

type FaceBookUser struct {
	ID    string `json:"id" valid:"required"`
	Email string `json:"email" valid:"required"`
	Name  string `json:"name" valid:"required"`
}

func (f *FaceBookUser) IsExist(ctx context.Context) []SocialGroup {
	var groups []SocialGroup
	err := factory.DB(ctx).
		Table("user").
		Join("INNER", "social_auth", "user.id = social_auth.user_id").
		Where("user.email = ?", f.Email).
		Find(&groups)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return groups
}

func (f *FaceBookUser) ToUser() *User {
	fmt.Println(f)
	u := User{
		Email:    f.Email,
		Nickname: f.Name,
	}
	return &u
}

func (f *FaceBookUser) ToSocial(id int64, token string) *SocialAuth {
	s := SocialAuth{
		SocialID:   f.ID,
		UserID:     id,
		SocialType: "FACEBOOK",
		Token:      token,
	}
	return &s
}
