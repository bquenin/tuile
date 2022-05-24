# Web Assembly Sample

## Prerequisites

 * Make sure the wasm Go helper is up to date:
```
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

 * Install [Caddy](https://caddyserver.com/docs/install)

## Build the sample

```
GOARCH=wasm GOOS=js go build -o main.wasm main.go
```

## Run the sample

```
caddy run
```
 * Go to http://localhost/
