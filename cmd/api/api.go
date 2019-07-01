package main

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"os/signal"
	"runtime"
	"showper_server/database"
	_ "showper_server/docs"
	"showper_server/middlewares"
	"syscall"
	"time"
)

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}

type App struct {
	*echo.Echo
}

func NewApp() *App {
	return &App{
		Echo: echo.New(),
	}
}

// @title Sumwhere API
// @version 2.0
// @description This is a Sumwhere server API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://www.sumwhere.kr
// @contact.email qjadn0914@naver.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host api.sumwhere.kr
// @BasePath /v1
// @schemes https http
func (a *App) Run(port string) error {

	a.importControllers()

	db, err := database.NewDatabase()
	if err != nil {
		return err
	}

	a.Static("/images", "/images")
	a.Use(middlewares.ContextDB("sumwhere", db))
	a.Use(middlewares.Logger())
	a.Pre(middleware.RemoveTrailingSlash())
	a.Use(middleware.CORS())
	a.Use(middleware.RequestID())
	a.Use(middleware.Recover())
	a.Use(middleware.JWT([]byte("bumwoopark")))

	a.Validator = &Validator{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := a.Shutdown(ctx); err != nil {
			a.Logger.Fatal(err)
		}
	}()
	return a.Start(fmt.Sprintf(":%s", port))
}

func (a *App) importControllers() {
	v1 := a.Group("/v1")
	v1.GET("/swagger/*", echoSwagger.WrapHandler)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	application := NewApp()
	if err := application.Run("8080"); err != nil {
		panic(err)
	}
}