package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	app "github.com/murakamiii/go-check-app-release/internal/app"
	slack "github.com/murakamiii/go-check-app-release/internal/slack"
	tag "github.com/murakamiii/go-check-app-release/internal/tag"
)

func main() {
	flag.Parse()
	slackPath := flag.Arg(0)

	app := app.App{&http.Client{Timeout: time.Second * 10}}
	v := map[string]string{}

	iosv, err := app.GetiOSVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	v["ios"] = iosv

	androidv, err := app.GetAndroidVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	v["android"] = androidv

	msgs, err := tag.UpdateVersionTags(v)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(slackPath) > 0 && len(msgs) > 0 {
		sl := slack.Slack{&http.Client{Timeout: time.Second * 10}}
		err := sl.PostMessages(slackPath, msgs)
		if err != nil {
			fmt.Println(err)
		}
	}
}
