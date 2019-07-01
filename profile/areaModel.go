package profile

import "time"

type Area struct {
	ID        int       `json:"id pk autoincr" xorm:"id"`
	City      string    `json:"city" xorm:"city"`
	District  string    `json:"district" xml:"district"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`
}
