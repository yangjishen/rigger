package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"rigger/rigger"
)

func main() {

	app := cli.NewApp()
	app.Version = "1.0.0-rc"
	app.Usage = "Generate scaffold project layout for Go."

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Generate scaffold project layout",
			Subcommands: []cli.Command{
				{
					Name:  "api",
					Usage: "Generate scaffold api project layout",
					Action: func(c *cli.Context) error {

						currDir, err := filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
						if err != nil {
							return err
						}
						fmt.Println(currDir)
						err = rigger.New(false).Generate(currDir)
						if err == nil {
							fmt.Println("Success Created. Please excute `make up` to start service.")
						}
						return err
					},
				},
				{
					Name:  "srv",
					Usage: "Generate scaffold srv project layout",
					Action: func(c *cli.Context) error {
						fmt.Println("Generate scaffold srv project layout", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
