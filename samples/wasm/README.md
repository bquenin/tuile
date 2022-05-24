# Web Assembly Sample

## Prerequisites

 * Make sure the wasm Go helper is up to date:
```
 $ cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

 * Install [Caddy](https://caddyserver.com/):
```
 $ brew install caddy
```

## Build the sample

```
 $ GOARCH=wasm GOOS=js go build -o main.wasm main.go
```

## Run the sample
```
$ caddy start
```
 * Go to http://localhost/
