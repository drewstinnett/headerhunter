package headerhunter

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, err := New()
	require.EqualError(t, err, "must set staticDir or proxyURL")

	_, err = New(WithStaticDir("foo"), WithProxyURL("bar"))
	require.EqualError(t, err, "only set either staticDir or proxyURL")
}

func TestStaticMode(t *testing.T) {
	b := bytes.NewBufferString("")
	h, err := New(
		WithWriter(b),
		WithStaticDir("./testdata/static"),
	)
	require.NoError(t, err)

	srv := httptest.NewServer(h)
	defer srv.Close()

	req := mustReq("GET", srv.URL+"/index.html", nil)
	req.Header.Set("foo", "bar")
	req.Header.Set("user-agent", "test")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	got, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "Hello fellow Hunter!\n", string(got))
	require.Contains(t, b.String(), `{"Accept-Encoding":["gzip"],"Foo":["bar"],"User-Agent":["test"]}`)
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
	require.NoError(t, err)

	srv := httptest.NewServer(h)
	defer srv.Close()

	req := mustReq("GET", srv.URL+"/index.html", nil)
	req.Header.Set("foo", "bar")
	req.Header.Set("user-agent", "test")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	got, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "Hello from upstream", string(got))
	require.Contains(t, b.String(), `{"Accept-Encoding":["gzip"],"Foo":["bar"],"User-Agent":["test"]}`)
}

func mustReq(method, url string, body io.Reader) *http.Request {
	got, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	return got
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
	require.NoError(t, err)

	srv := httptest.NewServer(h)
	defer srv.Close()

	req := mustReq("GET", srv.URL+"/prefix/index.html", nil)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	got, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "Hello from upstream", string(got))
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
	require.NoError(t, err)

	srv := httptest.NewServer(h)
	defer srv.Close()

	req := mustReq("POST", srv.URL+"/", bytes.NewReader([]byte("hi")))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	got, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "Hello from upstream", string(got))
	require.Contains(t, b.String(), `"size":2`)  // size of request
	require.Contains(t, b.String(), `"size":19`) // size of response
}
