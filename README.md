# mtls.dev

## What is it?
It compiles Cloudflare's [cfssl](https://github.com/cloudflare/cfssl) down to [WebAssembly](https://github.com/golang/go/wiki/WebAssembly) to generate mTLS certs in the browser.

## Why?
Generating mTLS certs and using them is more challenging than it needs to be.

## Developing
* `make wasm` - builds the wasm binary that gets loaded in the browser.
* `make dev-serve` - runs a simple python webserver for local dev, serving on port 8080.
* `make test` - runs tests.
* `make integration-build` and `make integration` - these test the sample code snippets, first serving and then using the client to hit the server.
