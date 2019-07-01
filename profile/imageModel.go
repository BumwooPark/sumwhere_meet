package profile

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sumwhere_meet/factory"
	"sumwhere_meet/utils"
	"time"
)

type Image struct {
	ID        int64     `json:"id" xorm:"id pk autoincr"`
	UserID    int64     `json:"userID" xorm:"user_id"`
	ImageURL  string    `json:"url" xorm:"image_url"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`
}

func (i Image) Save(ctx context.Context, UserID int64, file []*multipart.FileHeader) error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	dataPath := fmt.Sprintf("/images/%d/profile/", UserID)

	path := dir + dataPath
	utils.CreateDirIfNotExist(path)

	images := []Image{}

	for _, v := range file {

		filepath := path + v.Filename
		datafilePath := dataPath + v.Filename

		src, err := v.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		img := Image{
			UserID:   UserID,
			ImageURL: datafilePath,
		}

		images = append(images, img)
	}

	if err := i.save(ctx, images); err != nil {
		return err
	}

	return nil
}

func (Image) save(ctx context.Context, images []Image) error {
	_, err := factory.DB(ctx).Insert(images)
	if err != nil {
		return err
	}
	return nil
}
