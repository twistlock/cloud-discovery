package main

import (
	"flag"
	"fmt"
	"github.com/twistlock/cloud-discovery/internal/provider/aws"
	"os"
	"text/tabwriter"
)

var (
	username, password string
)

func main() {
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	flag.Parse()
	if username == "" {
		panic("username is missing")
	}
	if password == "" {
		panic("password is missing")
	}
	results := aws.Discover(username, password)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Type\tRegion\tID")
	for _, r := range results.Results {
		for _, asset := range r.Assets {
			fmt.Fprintf(w, "%s\t%s\t%s\n", r.Type, r.Region, asset.ID)
		}

	}
	w.Flush()
}
