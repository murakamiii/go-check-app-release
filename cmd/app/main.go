package main

import (
	"fmt"
	ios "github.com/murakamiii/go-check-app-release/internal/app/ios"
	"net/http"
	"time"
)

func main() {
	app := ios.App { &http.Client{ Timeout: time.Second * 10 } }
	version, err := app.GetVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(version)
}