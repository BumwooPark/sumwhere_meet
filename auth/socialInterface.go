package auth

import "context"

type SocialAction interface {
	IsExist(ctx context.Context) []SocialGroup
	ToUser() *User
	ToSocial(id int64, token string) *SocialAuth
}
