package httpserver

import (
	"bytes"
	"net/http"
	"testing"
)

func TestJSON(t *testing.T) {
	recorder := &responseRecorder{header: make(http.Header)}
	JSON(recorder, http.StatusOK, map[string]string{"status": "ok"})

	if recorder.code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected JSON content type, got %q", got)
	}
	if got := recorder.body.String(); got != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected body %q", got)
	}
}

type responseRecorder struct {
	header http.Header
	code   int
	body   bytes.Buffer
}

func (r *responseRecorder) Header() http.Header {
	return r.header
}

func (r *responseRecorder) WriteHeader(code int) {
	r.code = code
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	if r.code == 0 {
		r.code = http.StatusOK
	}
	return r.body.Write(data)
}
