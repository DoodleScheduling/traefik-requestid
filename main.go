package traefik_requestid

import (
	"context"
	"fmt"
	"net/http"
)

const defaultUpstreamHeaderName = "x-request-id"
const defaultDownstreamHeaderName = "x-request-id"

type Config struct {
	UpstreamHeaderName   string `json:"upstreamHeaderName"`
	DownstreamHeaderName string `json:"downstreamHeaderName"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		UpstreamHeaderName:   defaultUpstreamHeaderName,
		DownstreamHeaderName: defaultDownstreamHeaderName,
	}
}

// RequestIDHeader header if it's missing
type RequestIDHeader struct {
	UpstreamHeaderName   string
	DownstreamHeaderName string
	name                 string
	next                 http.Handler
}

// New created a new RequestIDHeader plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	tIDHdr := &RequestIDHeader{
		next: next,
		name: name,
	}

	if config == nil {
		return nil, fmt.Errorf("config can not be nil")
	}

	if config.UpstreamHeaderName == "" {
		tIDHdr.UpstreamHeaderName = defaultUpstreamHeaderName
	} else {
		tIDHdr.UpstreamHeaderName = config.UpstreamHeaderName
	}

	if config.DownstreamHeaderName == "" {
		tIDHdr.DownstreamHeaderName = defaultDownstreamHeaderName
	} else {
		tIDHdr.DownstreamHeaderName = config.DownstreamHeaderName
	}

	return tIDHdr, nil
}

func (t *RequestIDHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	randomUUID := newUUID().String()
	req.Header.Add(t.UpstreamHeaderName, randomUUID)
	rw.Header().Add(t.DownstreamHeaderName, randomUUID)
	t.next.ServeHTTP(rw, req)
}
