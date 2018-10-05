package main

import (
	"flag"
	"fmt"
	"github.com/twistlock/cloud-discovery/internal/provider/aws"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"os"
	"text/tabwriter"
)

func main() {
	var (
		username, password string
	)
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	flag.Parse()
	if username == "" {
		panic("username is missing")
	}
	if password == "" {
		panic("password is missing")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Type\tRegion\tID")
	aws.Discover(username, password, func(result shared.CloudDiscoveryResult) {
		for _, asset := range result.Assets {
			fmt.Fprintf(w, "%s\t%s\t%s\n", result.Type, result.Region, asset.ID)
		}
		w.Flush()
	})
}
