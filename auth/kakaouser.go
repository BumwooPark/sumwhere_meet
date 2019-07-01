package auth

import (
	"context"
	"fmt"
	"showper_server/factory"
	"strconv"
)

type KakaoUser struct {
	ID         int `json:"id"`
	Token      string
	Properties struct {
		Nickname       string `json:"nickname"`
		ProfileImage   string `json:"profile_image"`
		ThumbnailImage string `json:"thumbnail_image"`
	} `json:"properties"`
	KakaoAccount struct {
		HasEmail        bool   `json:"has_email"`
		IsEmailValid    bool   `json:"is_email_valid"`
		IsEmailVerified bool   `json:"is_email_verified"`
		Email           string `json:"email"`
		HasAgeRange     bool   `json:"has_age_range"`
		HasBirthday     bool   `json:"has_birthday"`
		HasGender       bool   `json:"has_gender"`
	} `json:"kakao_account"`
}

type SocialGroup struct {
	User       `xorm:"extends"`
	SocialAuth `xorm:"extends"`
}

func (k *KakaoUser) IsExist(ctx context.Context) []SocialGroup {
	var groups []SocialGroup
	err := factory.DB(ctx).
		Table("user").
		Join("INNER", "social_auth", "user.id = social_auth.user_id").
		Where("user.email = ?", k.KakaoAccount.Email).
		Find(&groups)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return groups
}

func (k *KakaoUser) ToUser() *User {
	u := User{
		Email:    k.KakaoAccount.Email,
		Nickname: k.Properties.Nickname,
	}
	return &u
}

func (k *KakaoUser) ToSocial(id int64, token string) *SocialAuth {
	s := SocialAuth{
		SocialID:   strconv.Itoa(k.ID),
		UserID:     id,
		SocialType: "KAKAO",
		Token:      token,
	}
	return &s
}
