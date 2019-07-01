package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/testfixtures.v2"
	"log"
	"os"
	"testing"
)

var (
	db *sql.DB
	fixtures *testfixtures.Context
)

func TestMain(m *testing.M){
	db, err := sql.Open("mysql","root:1q2w3e4r@tcp(localhost:3306)/test?charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}

	fixtures, err = testfixtures.NewFiles(db,&testfixtures.MySQL{},"user.yml")
	if err != nil {
		log.Fatal(err)
	}
	testfixtures.SkipDatabaseNameCheck(true)

	if err := fixtures.DetectTestDatabase(); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		fmt.Println(err)
	}
}

func TestX(t *testing.T){
	prepareTestDatabase()
}