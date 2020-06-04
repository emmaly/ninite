package main

import (
	"fmt"
	"os"

	niniteclassic "github.com/emmaly/ninite/classic"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("An application to be uninstalled must be indicated.")
		os.Exit(1)
	}

	nc, err := niniteclassic.New(".")
	if err != nil {
		panic(err)
	}

	as := make(chan niniteclassic.Status)
	if err := nc.Select(os.Args[1:]...).Uninstall(as); err != nil {
		panic(err)
	}

	for app := range as {
		if app.Error != nil {
			panic(app.Error)
		}

		fmt.Printf("[%s]\n\tStatus: %s\n\tReason: %s\n\n", app.App, app.Status, app.Reason)
	}
}
