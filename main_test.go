package traefik_requestid

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var uuidV4Re = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestCreateConfig(t *testing.T) {
	cfg := CreateConfig()
	assert.Equal(t, defaultUpstreamHeaderName, cfg.UpstreamHeaderName)
	assert.Equal(t, defaultDownstreamHeaderName, cfg.DownstreamHeaderName)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name                 string
		cfg                  *Config
		errExpected          bool
		upstreamHeaderName   string
		downstreamHeaderName string
	}{
		{
			name:        "nil config returns error",
			cfg:         nil,
			errExpected: true,
		},
		{
			name:                 "empty header names use defaults",
			cfg:                  &Config{},
			upstreamHeaderName:   defaultUpstreamHeaderName,
			downstreamHeaderName: defaultDownstreamHeaderName,
		},
		{
			name: "custom header names are used",
			cfg: &Config{
				UpstreamHeaderName:   "x-custom-upstream",
				DownstreamHeaderName: "x-custom-downstream",
			},
			upstreamHeaderName:   "x-custom-upstream",
			downstreamHeaderName: "x-custom-downstream",
		},
		{
			name: "default config values",
			cfg:  CreateConfig(),

			upstreamHeaderName:   defaultUpstreamHeaderName,
			downstreamHeaderName: defaultDownstreamHeaderName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := New(t.Context(), http.NotFoundHandler(), tt.cfg, tt.name)
			if tt.errExpected {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, handler)

			h, ok := handler.(*RequestIDHeader)
			require.True(t, ok)
			assert.Equal(t, tt.upstreamHeaderName, h.UpstreamHeaderName)
			assert.Equal(t, tt.downstreamHeaderName, h.DownstreamHeaderName)
		})
	}
}

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		name                 string
		cfg                  *Config
		upstreamHeaderName   string
		downstreamHeaderName string
	}{
		{
			name:                 "default headers are set on request and response",
			cfg:                  CreateConfig(),
			upstreamHeaderName:   defaultUpstreamHeaderName,
			downstreamHeaderName: defaultDownstreamHeaderName,
		},
		{
			name: "custom upstream and downstream headers are set",
			cfg: &Config{
				UpstreamHeaderName:   "x-custom-upstream",
				DownstreamHeaderName: "x-custom-downstream",
			},
			upstreamHeaderName:   "x-custom-upstream",
			downstreamHeaderName: "x-custom-downstream",
		},
		{
			name: "same header name for upstream and downstream",
			cfg: &Config{
				UpstreamHeaderName:   "x-req-id",
				DownstreamHeaderName: "x-req-id",
			},
			upstreamHeaderName:   "x-req-id",
			downstreamHeaderName: "x-req-id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedRequestHeader string
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				capturedRequestHeader = req.Header.Get(tt.upstreamHeaderName)
			})

			handler := &RequestIDHeader{
				next:                 next,
				name:                 tt.name,
				UpstreamHeaderName:   tt.cfg.UpstreamHeaderName,
				DownstreamHeaderName: tt.cfg.DownstreamHeaderName,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rw := httptest.NewRecorder()
			handler.ServeHTTP(rw, req)

			assert.Regexp(t, uuidV4Re, capturedRequestHeader, "upstream request header should be a valid UUID v4")

			responseHeader := rw.Header().Get(tt.downstreamHeaderName)
			assert.Regexp(t, uuidV4Re, responseHeader, "downstream response header should be a valid UUID v4")

			assert.Equal(t, capturedRequestHeader, responseHeader, "upstream and downstream should carry the same UUID")
		})
	}
}

func TestServeHTTP_NextHandlerIsCalled(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		called = true
		rw.WriteHeader(http.StatusAccepted)
	})

	handler := &RequestIDHeader{
		next:                 next,
		name:                 "test",
		UpstreamHeaderName:   defaultUpstreamHeaderName,
		DownstreamHeaderName: defaultDownstreamHeaderName,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)

	assert.True(t, called, "next handler must be called")
	assert.Equal(t, http.StatusAccepted, rw.Code)
}

func TestServeHTTP_ExistingHeaderIsNotOverwritten(t *testing.T) {
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler := &RequestIDHeader{
		next:                 next,
		name:                 "test",
		UpstreamHeaderName:   defaultUpstreamHeaderName,
		DownstreamHeaderName: defaultDownstreamHeaderName,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(defaultUpstreamHeaderName, "existing-id")
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)

	values := req.Header[http.CanonicalHeaderKey(defaultUpstreamHeaderName)]
	assert.Len(t, values, 2, "Add should append a new value, not replace the existing one")
	assert.Equal(t, "existing-id", values[0])
	assert.Regexp(t, uuidV4Re, values[1])
}
