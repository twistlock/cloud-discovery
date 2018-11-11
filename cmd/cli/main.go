package main

import (
	"encoding/json"
	"fmt"
	"github.com/twistlock/cloud-discovery/internal/nmap"
	"github.com/twistlock/cloud-discovery/internal/provider"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "cloud-discovery"
	app.Usage = " Cloud Discovery provides a point in time enumeration of all the cloud native platform services"
	app.Version = "1.0.0"

	var configPath, format, subnet string
	app.Commands = []cli.Command{
		{
			Name:  "discover",
			Usage: "Discover all cloud assets",
			Flags: []cli.Flag{cli.StringFlag{
				Name:        "config",
				Usage:       "Path to credential configuration",
				Destination: &configPath,
			},
				cli.StringFlag{
					Name:        "format",
					Usage:       "Output Formatting (json or csv)",
					Value:       "csv",
					Destination: &format,
				},
			},
			Action: func(c *cli.Context) error {
				if configPath == "" {
					return fmt.Errorf("missing config path")
				}
				data, err := ioutil.ReadFile(configPath)
				if err != nil {
					return err
				}
				var creds []shared.Credentials
				if err := json.Unmarshal(data, &creds); err != nil {
					return err
				}
				provider.Discover(creds, os.Stdout, shared.Format(format))
				return nil
			},
		},
		{
			Name:  "nmap",
			Usage: "Scan all exposed cloud assets",
			Flags: []cli.Flag{cli.StringFlag{
				Name:        "subnet",
				Usage:       "The subnet to scan",
				Value:       "127.0.0.1",
				Destination: &subnet,
			},
			},
			Action: func(c *cli.Context) error {
				nmap.Nmap(os.Stdout, subnet, true)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
