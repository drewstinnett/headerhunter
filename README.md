<p align="center">
  <img alt="GoFlex Logo" src="./static/images/logo.png" height="200" />
  <h3 align="center">HeaderHunter</h3>
  <p align="center">Hunt down headers!</p>
</p>

# HeaderHunter

A lightweight HTTP debugging tool that acts as a reverse proxy or static file server, capturing and logging all HTTP request and response headers in JSON format. Perfect for debugging web applications, inspecting API communications, and understanding HTTP protocol details.

## Features

- **Reverse Proxy Mode**: Proxy requests to any HTTP/HTTPS website
- **Static File Server Mode**: Serve files from a local directory
- **JSON Logging**: Structured JSON output for easy parsing and analysis
- **TLS/HTTPS Support**: Run with custom TLS certificates
- **Configurable**: Customizable listen address, port, timeouts, and URL prefixes
- **Request/Response Tracking**: Captures headers, methods, URLs, status codes, and body sizes
- **Clean Shutdown**: Graceful shutdown with configurable timeouts

## Installation

### Using Go Install

```bash
go install github.com/drewstinnett/headerhunter/cmd/headerhunter@latest
```

### Building from Source

```bash
git clone https://github.com/drewstinnett/headerhunter.git
cd headerhunter
go build -o headerhunter ./cmd/headerhunter
```

### Requirements

- Go 1.23 or higher

## Usage

### Basic Examples

**Proxy to a website:**
```bash
headerhunter serve https://www.google.com
```

**Serve static files:**
```bash
headerhunter serve ./public
```

**Custom port:**
```bash
headerhunter serve https://api.example.com -a :8080
```

**With TLS:**
```bash
headerhunter serve https://example.com --cert server.crt --key server.key
```

**With URL prefix:**
```bash
headerhunter serve https://example.com -p /api/
```

### Making Requests

Once running, make requests to the local server:

```bash
curl http://127.0.0.1:3000/
```

### Example Output

```bash
$ headerhunter serve https://www.google.com
2025/02/01 16:42:46 INFO headerhunter 🫨 : launching http server addr=127.0.0.1:3000
{
  "type": "request",
  "time": "2025-04-12T17:56:24.573701-04:00",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "User-Agent": [
      "curl/8.7.1"
    ]
  },
  "remote_addr": "[::1]:64981",
  "method": "GET",
  "url": "/",
  "size": 0
}
{
  "type": "response",
  "time": "2025-04-12T17:56:24.871269-04:00",
  "status_code": 200,
  "headers": {
    "Accept-Ch": [
      "Sec-CH-Prefers-Color-Scheme"
    ],
    "Cache-Control": [
      "private, max-age=0"
    ],
    "Content-Type": [
      "text/html; charset=ISO-8859-1"
    ],
    "Date": [
      "Sat, 12 Apr 2025 21:56:24 GMT"
    ],
    "Server": [
      "gws"
    ],
    "Set-Cookie": [
      "AEC=AVcja2dDI6FmBPH-wqB_aEF5579VMSn6N99mOhSIUINS8rWJ3rQkt1LQbQ; expires=Thu, 09-Oct-2025 21:56:24 GMT; path=/; domain=.google.com; Secure; HttpOnly; SameSite=lax"
    ],
    "X-Frame-Options": [
      "SAMEORIGIN"
    ]
  },
  "method": "GET",
  "url": "/",
  "size": 17102
}
```

## Command-Line Options

```
Usage:
  headerhunter serve DIR|URL [flags]

Flags:
  -a, --addr string              Listen address (default ":3000")
  -c, --cert string              TLS certificate file
  -k, --key string               TLS key file (required with --cert)
  -p, --prefix string            URL prefix for routing (default "/")
      --read-timeout duration    Server read timeout (default 10m)
      --write-timeout duration   Server write timeout (default 10m)
  -v, --verbose                  Enable verbose logging
  -h, --help                     Help for serve
```

## Output Format

HeaderHunter outputs two types of JSON logs:

### Request Log

```json
{
  "type": "request",
  "time": "2025-04-12T17:56:24.573701-04:00",
  "headers": {
    "Accept": ["*/*"],
    "User-Agent": ["curl/8.7.1"]
  },
  "remote_addr": "[::1]:64981",
  "method": "GET",
  "url": "/",
  "size": 0
}
```

### Response Log

```json
{
  "type": "response",
  "time": "2025-04-12T17:56:24.871269-04:00",
  "status_code": 200,
  "headers": {
    "Content-Type": ["text/html; charset=UTF-8"]
  },
  "method": "GET",
  "url": "/",
  "size": 17102
}
```

## Use Cases

- **Debugging API Integrations**: See exactly what headers are being sent and received
- **Testing Reverse Proxy Configurations**: Verify proxy behavior and header forwarding
- **Learning HTTP Protocol**: Understand how browsers and servers communicate
- **Cookie/Session Debugging**: Inspect Set-Cookie headers and session management
- **Security Analysis**: Examine security headers (CORS, CSP, HSTS, etc.)
- **Middleware Development**: Test custom middleware header modifications
- **Performance Analysis**: Track response times and content sizes

## Troubleshooting

**Port already in use:**
```bash
# Use a different port
headerhunter serve https://example.com -a :8080
```

**Certificate errors with HTTPS proxying:**
HeaderHunter proxies HTTPS sites over HTTP by default. To serve over HTTPS, use the `--cert` and `--key` flags.

**Large response bodies:**
HeaderHunter logs body sizes but not content. Use verbose mode (`-v`) for more detailed logging.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See LICENSE file for details
