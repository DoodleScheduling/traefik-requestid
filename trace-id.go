package traefik_add_trace_id

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

const defaultTraceID = "X-Trace-Id"

// Config the plugin configuration.
type Config struct {
	HeaderPrefix string `json:"headerPrefix"`
	HeaderName   string `json:"headerName"`
	Verbose      bool   `json:"verbose"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderPrefix: "",
		HeaderName:   defaultTraceID,
	}
}

// TraceIDHeader header if it's missing
type TraceIDHeader struct {
	headerName   string
	headerPrefix string
	name         string
	next         http.Handler
	verbose      bool
}

// New created a new TraceIDHeader plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config == nil {
		return nil, fmt.Errorf("config can not be nil")
	}

	tIDHdr := &TraceIDHeader{
		next:         next,
		name:         name,
		verbose:      config.Verbose,
		headerPrefix: config.HeaderPrefix,
	}

	if config.HeaderName == "" {
		tIDHdr.headerName = defaultTraceID
	} else {
		tIDHdr.headerName = config.HeaderName
	}

	return tIDHdr, nil
}

func (t *TraceIDHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	traceID := req.Header.Get(t.headerName)

	if traceID == "" {
		traceID = fmt.Sprintf("%s%s", t.headerPrefix, newUUID().String())
		req.Header.Set(t.headerName, traceID)
	}

	if t.verbose {
		log.Println(traceID)
	}

	rw.Header().Set(t.headerName, traceID)

	t.next.ServeHTTP(rw, req)
}
