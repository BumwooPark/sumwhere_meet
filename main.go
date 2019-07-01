package main

import (
	"runtime"
	"showper_server/app"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	application := app.NewApp()
	if err := application.Run("8080"); err != nil {
		panic(err)
	}
}
