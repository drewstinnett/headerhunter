package headerhunter

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		args    []Option
		wantErr string
	}{
		{args: nil, wantErr: "must set staticDir or proxyURL"},
		{args: []Option{WithStaticDir("foo"), WithProxyURL("bar")}, wantErr: "only set either staticDir or proxyURL"},
	}

	for _, tt := range tests {
		_, err := New(tt.args...)
		if err == nil || err.Error() != tt.wantErr {
			t.Errorf("expected error %q, got %v", tt.wantErr, err)
		}
	}
}

func TestStaticMode(t *testing.T) {
	b := bytes.NewBufferString("")
	h, err := New(
		WithWriter(b),
		WithStaticDir("./testdata/static"),
	)
	if err != nil {
		t.Fatalf("failed to create new Hunter: %v", err)
	}

	srv := httptest.NewServer(h)
	defer srv.Close()

	req := mustReq(context.Background(), "GET", srv.URL+"/index.html", nil)
	req.Header.Set("Foo", "bar")
	req.Header.Set("User-Agent", "test")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("failed to make request: %v", err)
	}
	defer dclose(resp.Body)

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}
	if got, expected := string(got), "Hello fellow Hunter!\n"; got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}

	if got, expected := b.String(), `{"Accept-Encoding":["gzip"],"Foo":["bar"],"User-Agent":["test"]}`; !bytes.Contains([]byte(got), []byte(expected)) {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestProxyMode(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("Hello from upstream")); err != nil {
			panic(err)
		}
	}))
	defer upstream.Close()

	b := bytes.NewBufferString("")
	h, err := New(
		WithWriter(b),
		WithProxyURL(upstream.URL),
	)
	if err != nil {
		t.Fatalf("failed to create new Hunter: %v", err)
	}

	srv := httptest.NewServer(h)
	defer srv.Close()

	req := mustReq(context.Background(), "GET", srv.URL+"/index.html", nil)
	req.Header.Set("Foo", "bar")
	req.Header.Set("User-Agent", "test")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer dclose(resp.Body)

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if got, expect := string(got), "Hello from upstream"; got != expect {
		t.Fatalf("expected %q, got %q", expect, got)
	}
	if got, expect := b.String(), `{"Accept-Encoding":["gzip"],"Foo":["bar"],"User-Agent":["test"]}`; !bytes.Contains([]byte(got), []byte(expect)) {
		t.Fatalf("expected %q, got %q", expect, got)
	}
}

func TestPrefix(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("Hello from upstream")); err != nil {
			panic(err)
		}
	}))
	defer upstream.Close()

	h, err := New(
		WithPrefix("/prefix"),
		WithProxyURL(upstream.URL),
	)
	if err != nil {
		t.Fatalf("failed to create new Hunter: %v", err)
	}

	srv := httptest.NewServer(h)
	defer srv.Close()

	req := mustReq(context.Background(), "GET", srv.URL+"/prefix/index.html", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer dclose(resp.Body)

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if got, expect := string(got), "Hello from upstream"; got != expect {
		t.Fatalf("expected %q, got %q", expect, got)
	}
}

func TestSizeReporting(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("Hello from upstream")); err != nil {
			panic(err)
		}
	}))
	defer upstream.Close()

	b := bytes.NewBufferString("")
	h, err := New(
		WithWriter(b),
		WithProxyURL(upstream.URL),
	)
	if err != nil {
		t.Fatalf("failed to create new Hunter: %v", err)
	}

	srv := httptest.NewServer(h)
	defer srv.Close()

	req := mustReq(context.Background(), "POST", srv.URL+"/", bytes.NewReader([]byte("hi")))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer dclose(resp.Body)

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if got, expect := string(got), "Hello from upstream"; got != expect {
		t.Fatalf("expected %q, got %q", expect, got)
	}
	if got, expect := b.String(), `"size":2`; !bytes.Contains([]byte(got), []byte(expect)) {
		t.Fatalf("expected %q, got %q", expect, got)
	}
	if got, expect := b.String(), `"size":19`; !bytes.Contains([]byte(got), []byte(expect)) {
		t.Fatalf("expected %q, got %q", expect, got)
	}
}

func dclose(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}

func mustReq(ctx context.Context, method, url string, body io.Reader) *http.Request {
	got, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		panic(err)
	}
	return got
}
