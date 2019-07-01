package profile

import "time"

type InterestCategory struct {
	ID        int       `json:"id" xorm:"id pk autoincr"`
	Name      string    `json:"name" xorm:"name VARCHAR(20)"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`
}

type InterestDetail struct {
	ID         int       `json:"id" json:"id pk autoincr"`
	CategoryID int       `json:"categoryID" xorm:"category_id"`
	Name       string    `json:"name" xorm:"name VARCHAR(30)"`
	CreatedAt  time.Time `json:"createdAt" xorm:"created"`
	UpdateAt   time.Time `json:"updateAt"xorm:"updated"`
}
