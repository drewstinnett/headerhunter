/*
Package headerhunter proxies web requests and logs all of the headers
*/
package headerhunter

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

// Hunter contains all of the methods and such to do the headerhunter
type Hunter struct {
	writer io.Writer
	// mode      Mode
	handler        http.Handler
	prefix         string
	staticDir      string
	proxyURLString string
	proxyURL       *url.URL
}

type requestLog struct {
	Type       string              `json:"type"`
	Time       time.Time           `json:"time"`
	Headers    map[string][]string `json:"headers"`
	RemoteAddr string              `json:"remote_addr"`
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Size       int                 `json:"size"`
}
type responseLog struct {
	Type       string              `json:"type"`
	Time       time.Time           `json:"time"`
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Size       int                 `json:"size"`
}

func headerMap(h http.Header) map[string][]string {
	m := map[string][]string{}
	for k, v := range h {
		m[k] = append([]string{}, v...)
	}
	return m
}

// Handler returns an http.Handler based on the parameters defined in the Hunter
func (h Hunter) Handler() (http.Handler, error) {
	mux := http.NewServeMux()
	var handler http.Handler
	switch {
	case h.staticDir != "":
		handler = http.FileServer(http.Dir(h.staticDir))
	case h.proxyURL != nil:
		handler = httputil.NewSingleHostReverseProxy(h.proxyURL)
	default:
		return nil, errors.New("unknown mode, no staticDir or remoteURL specified")
	}
	mux.Handle(h.prefix, handler)
	return mux, nil
}

// ServeHTTP serves up the Hunter
func (h Hunter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.logRequest(r); err != nil {
		slog.Warn("error logging request", "error", err)
	}

	rl := &responseLogger{ResponseWriter: w, statusCode: http.StatusOK}

	h.handler.ServeHTTP(rl, r)

	if err := h.logResponse(rl, r); err != nil {
		slog.Warn("error logging response", "error", err)
	}
}

// Option is passed in to New to define it's behavior
type Option func(*Hunter)

// New returns a new Hunter object using functional options
func New(opts ...Option) (*Hunter, error) {
	h := &Hunter{
		writer: os.Stdout,
		prefix: "/",
	}
	for _, opt := range opts {
		opt(h)
	}

	if h.staticDir == "" && h.proxyURLString == "" {
		return nil, errors.New("must set staticDir or proxyURL")
	}

	if h.staticDir != "" && h.proxyURLString != "" {
		return nil, errors.New("only set either staticDir or proxyURL")
	}

	if h.proxyURLString != "" {
		u, err := url.Parse(h.proxyURLString)
		if err != nil {
			return nil, err
		}
		h.proxyURL = u
	}

	handler, err := h.Handler()
	if err != nil {
		return nil, err
	}
	h.handler = handler
	return h, nil
}

// WithWriter sets the writer for a Hunter
func WithWriter(w io.Writer) Option {
	return func(h *Hunter) {
		h.writer = w
	}
}

// WithStaticDir sets the static directory for a StaticMode functionality
func WithStaticDir(s string) Option {
	return func(h *Hunter) {
		h.staticDir = s
	}
}

// WithProxyURL sets the ProxyURL to proxy all the stuff
func WithProxyURL(s string) Option {
	return func(h *Hunter) {
		h.proxyURLString = s
	}
}

// WithPrefix sets the URL prefix for all the forwarding
func WithPrefix(s string) Option {
	return func(h *Hunter) {
		// Always ensure the string ends in a single /
		h.prefix = strings.TrimSuffix(s, "/") + "/"
	}
}

type responseLogger struct {
	http.ResponseWriter
	statusCode int
	headers    http.Header
	bodySize   int
}

func (rl *responseLogger) WriteHeader(code int) {
	rl.statusCode = code
	rl.headers = rl.ResponseWriter.Header()
	rl.ResponseWriter.WriteHeader(code)
}

func (rl *responseLogger) Write(b []byte) (int, error) {
	n, err := rl.ResponseWriter.Write(b)
	rl.bodySize += n // Track response body size
	return n, err
}

func (h Hunter) logResponse(rl *responseLogger, r *http.Request) error {
	out, err := json.Marshal(responseLog{
		Type:       "response",
		Time:       time.Now(),
		StatusCode: rl.statusCode,
		Headers:    headerMap(rl.headers),
		Method:     r.Method,
		URL:        r.URL.String(),
		Size:       rl.bodySize,
	})
	if err != nil {
		return err
	}
	_, err = h.writer.Write(out)
	if err != nil {
		return err
	}
	_, err = h.writer.Write([]byte("\n"))
	return err
}

func readBodySize(r *http.Request) (int, error) {
	var bodySize int
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return 0, err
		}
		bodySize = len(bodyBytes)
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes))) // Restore request body
	}
	return bodySize, nil
}

func (h Hunter) logRequest(r *http.Request) error {
	// Read and measure request body
	bs, err := readBodySize(r)
	if err != nil {
		return err
	}

	out, err := json.Marshal(requestLog{
		Type:       "request",
		Time:       time.Now(),
		RemoteAddr: r.RemoteAddr,
		Headers:    headerMap(r.Header),
		Method:     r.Method,
		URL:        r.URL.String(),
		Size:       bs,
	})
	if err != nil {
		return err
	}
	out = append(out, []byte("\n")...)
	_, err = h.writer.Write(out)
	if err != nil {
		return err
	}
	return err
}
