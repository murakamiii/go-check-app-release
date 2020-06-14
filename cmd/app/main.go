package main

import (
	"fmt"
	ios "github.com/murakamiii/go-check-app-release/internal/app/ios"
)

func main() {
	version, err := ios.GetVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(version)
}