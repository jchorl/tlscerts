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

## Why the forks?
* [cfssl](https://github.com/cloudflare/cfssl)
  * cfssl doesn't build for wasm, a dependency needs to be upgraded ([PR](https://github.com/cloudflare/cfssl/pull/1059))
  * fork here: https://github.com/jchorl/cfssl
* golang docker image
  * [glog](https://github.com/golang/glog) crashes on init in the browser due to missing syscalls
  * Turns out, there's an [issue](https://github.com/golang/go/issues/34627) for this. And a [fix](https://go-review.googlesource.com/c/go/+/199698/).
  * As of writing, this fix hasn't been released in a golang release and docker hub doesn't host nightly builds of golang.
  * So I built a master golang image with:
    ```dockerfile
    FROM golang:1.13.4

    RUN git clone https://go.googlesource.com/go goroot && \
        cd goroot && \
        git checkout master && \
        cd src && \
        ./make.bash

    ENV PATH="/go/goroot/bin:${PATH}"
    ```
