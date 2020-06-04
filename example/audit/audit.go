package main

import (
	"flag"
	"fmt"

	niniteclassic "github.com/emmaly/ninite/classic"
)

func main() {
	showOnlyInstalled := flag.Bool("installed", false, "Show only installed apps")
	flag.Parse()

	nc, err := niniteclassic.New(".")
	if err != nil {
		panic(err)
	}

	ac := make(chan niniteclassic.AppAudit)
	if err := nc.Audit(ac); err != nil {
		panic(err)
	}

	for app := range ac {
		if *showOnlyInstalled && !app.Installed {
			continue
		}

		fmt.Printf("[%s]\n\tStatus: %s\n\tVersion: %s\n\tInstalled: %t\n\n", app.App, app.Status, app.Version, app.Installed)
	}
}
