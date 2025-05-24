<p align="center">
  <img alt="GoFlex Logo" src="./static/images/logo.png" height="200" />
  <h3 align="center">HeaderHunter</h3>
  <p align="center">Hunt down headers!</p>
</p>

# Header Hunter

Runs a simple reverse proxy against a website or static directory, and prints
the request and response headers

## Usage

This is what would be logged if using something like:

```bash
curl -X POST http://127.0.0.1:3000/index.html
...
```

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
    "Alt-Svc": [
      "h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000"
    ],
    "Cache-Control": [
      "private, max-age=0"
    ],
    "Content-Security-Policy-Report-Only": [
      "object-src 'none';base-uri 'self';script-src 'nonce-swri3VTQ0WsyjepNhVSTFg' 'strict-dynamic' 'report-sample' 'unsafe-eval' 'unsafe-inline' https: http:;report-uri https://csp.withgoogle.com/csp/gws/oth
er-hp"
    ],
    "Content-Type": [
      "text/html; charset=ISO-8859-1"
    ],
    "Date": [
      "Sat, 12 Apr 2025 21:56:24 GMT"
    ],
    "Expires": [
      "-1"
    ],
    "P3p": [
      "CP=\"This is not a P3P policy! See g.co/p3phelp for more info.\""
    ],
    "Server": [
      "gws"
    ],
    "Set-Cookie": [
      "AEC=AVcja2dDI6FmBPH-wqB_aEF5579VMSn6N99mOhSIUINS8rWJ3rQkt1LQbQ; expires=Thu, 09-Oct-2025 21:56:24 GMT; path=/; domain=.google.com; Secure; HttpOnly; SameSite=lax",
      "NID=523=RmarybA4EDHDXwWnIkXSyfhAhYVtEi96bFSiHFRrWwY2TWRLftNqtTvf-oeJtaadCayKD-usYH1OUp7cN_yNPU2JCsjLNkmzPSB_rN2zAIIQ2-IQ0gZhsk23-87V6vrdNh5Rvn7S17EPmY26DeReF0STsicxXoubll3VYGIQRRySI0gjfNaedPfgefKDsHgTK
kPUL3rKMOxav_UAxx5n76h1; expires=Sun, 12-Oct-2025 21:56:24 GMT; path=/; domain=.google.com; HttpOnly"
    ],
    "X-Frame-Options": [
      "SAMEORIGIN"
    ],
    "X-Xss-Protection": [
      "0"
    ]
  },
  "method": "GET",
  "url": "/",
  "size": 17102
}
```
