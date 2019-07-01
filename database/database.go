package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"os"
	"sumwhere_meet/utils"
)

func NewDatabase() (*xorm.Engine, error) {
	var url string
	dbUser := utils.DefaultEnv("DATABASE_USER", "root")
	database := utils.DefaultEnv("DATABASE_DRIVER", "mysql")
	dbPass := utils.DefaultEnv("DATABASE_PASS", "1q2w3e4r")
	dbName := utils.DefaultEnv("DATABASE_NAME", "test")

	switch os.Getenv("RELEASE_SYSTEM") {
	case "kubernetes":
		url = fmt.Sprintf("%s:%s@tcp(mysql-svc.sumwhere:3306)/%s", dbUser, dbPass, dbName)
	default:
		url = fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4", "root", "1q2w3e4r", dbName)
	}

	db, err := xorm.NewEngine(database, url)
	if err != nil {
		return nil, err
	}

	db.ShowSQL(true)
	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db, nil
}
