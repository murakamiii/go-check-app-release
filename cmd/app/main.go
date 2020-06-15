package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	ios "github.com/murakamiii/go-check-app-release/internal/app/ios"
	slack "github.com/murakamiii/go-check-app-release/internal/slack"
)

func main() {
	flag.Parse()
	slackPath := flag.Arg(0)

	app := ios.App{&http.Client{Timeout: time.Second * 10}}
	version, err := app.GetVersion()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(slackPath) > 0 {
		sl := slack.Slack{&http.Client{Timeout: time.Second * 10}}
		err := sl.PostMessage(slackPath, version)
		if err != nil {
			fmt.Println(err)
		}
	}
}
