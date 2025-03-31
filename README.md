# Header Hunter

Runs a simple reverse proxy against a website or static directory, and prints
the request and response headers

## Usage

```bash
$ headerhunter serve https://www.google.com --addr 127.0.0.1:3000
2025/02/01 16:42:46 INFO headerhunter 🫨 : launching http server addr=127.0.0.1:3000
## curl -X POST http://127.0.0.1:3000/index.html -d "foo"
{"type":"request","time":"2025-02-01T16:43:06.689213-05:00","headers":{"Accept":["*/*"],"Content-Length":["3"],"Content-Type":["application/x-www-form-urlencoded"],"User-Agent":["curl/8.7.1"]},"remote_addr":"
127.0.0.1:63538","method":"POST","url":"/index.html","size":3}
{"type":"response","time":"2025-02-01T16:43:06.88428-05:00","status_code":404,"headers":{"Alt-Svc":["h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000"],"Content-Length":["1571"],"Content-Type":["text/html;
charset=UTF-8"],"Date":["Sat, 01 Feb 2025 21:43:07 GMT"],"Referrer-Policy":["no-referrer"]},"method":"POST","url":"/index.html","size":1571}
```
