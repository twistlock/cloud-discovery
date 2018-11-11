package provider

import (
	"fmt"
	"github.com/twistlock/cloud-discovery/internal/provider/aws"
	"github.com/twistlock/cloud-discovery/internal/provider/gcp"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"io"
	"text/tabwriter"
)

func Discover(creds []shared.Credentials, wr io.Writer, format shared.Format) {
	var writer ResponseWriter
	if format == shared.FormatJson {
		writer = shared.NewJsonResponseWriter(wr)
	} else {
		writer = NewTabResponseWriter(wr)
	}
	for _, cred := range creds {
		switch cred.Provider {
		case shared.ProviderGCP:
			gcp.Discover(cred.Secret, writer.Write)
		default:
			aws.Discover(cred.ID, cred.Secret, writer.Write)
		}
	}
}

type ResponseWriter interface {
	Write(shared.CloudDiscoveryResult)
}

type csvResponseWriter struct {
	tw *tabwriter.Writer
}

func NewTabResponseWriter(writer io.Writer) *csvResponseWriter {
	tw := shared.NewTabWriter(writer)
	fmt.Fprintf(tw, "Type\tRegion\tID\n")
	return &csvResponseWriter{tw: tw}
}

func (w *csvResponseWriter) Write(result shared.CloudDiscoveryResult) {
	for _, asset := range result.Assets {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\n", result.Type, result.Region, asset.ID)
	}
	w.tw.Flush()
}
