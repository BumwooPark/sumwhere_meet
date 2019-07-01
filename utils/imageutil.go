package utils

import (
	"os"
)

//
//func ProfileSaver(header *multipart.FileHeader, user *models.User, imageName string) string {
//
//	if header == nil {
//		return ""
//	}
//
//	file, err := header.Open()
//	if err != nil {
//		return ""
//	}
//
//	defer file.Close()
//	path := fmt.Sprintf("/images/%d/profile/", user.Id)
//	CreateDirIfNotExist(path)
//
//	dst, err := os.Create(path + fmt.Sprintf("%s.jpg", imageName))
//	if err != nil {
//		return ""
//	}
//	defer dst.Close()
//	if _, err = io.Copy(dst, file); err != nil {
//		return ""
//	}
//	return fmt.Sprintf("/%d/profile/%s.jpg", user.Id, imageName)
//}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
