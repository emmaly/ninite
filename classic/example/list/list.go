package main

import (
	"flag"
	"fmt"

	niniteclassic "github.com/emmaly/ninite/classic"
)

func main() {
	showAlternate := flag.Bool("showAlternate", false, "Include alternate app versions")
	flag.Parse()

	nc, err := niniteclassic.New(".")
	if err != nil {
		panic(err)
	}

	av := make(chan niniteclassic.AppVersion)
	if err := nc.List(av); err != nil {
		panic(err)
	}

	for app := range av {
		if app.Error != nil {
			panic(app.Error)
		}

		if !*showAlternate && app.AlternateVersion {
			continue
		}

		fmt.Printf("[%s]\n\tVersion: %s\n\tCurrentVersion: %t\n\tAlternateVersion: %t\n\n", app.App, app.Version, app.CurrentVersion, app.AlternateVersion)
	}
}
