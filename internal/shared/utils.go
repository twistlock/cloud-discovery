package shared

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"text/tabwriter"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func NewTabWriter(wr io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(wr, 0, 0, 5, ' ', tabwriter.TabIndent)
}

type jsonResposeWriter struct {
	w io.Writer
}

func NewJsonResponseWriter(w io.Writer) *jsonResposeWriter {
	return &jsonResposeWriter{w: w}
}

func (w *jsonResposeWriter) Write(result CloudDiscoveryResult) {
	out, _ := json.Marshal(result)
	w.w.Write(out)
	w.w.Write([]byte("\n"))
}
