FROM golang:1.13.4

RUN git clone --depth 1 https://go.googlesource.com/go goroot && \
    cd goroot && \
    git checkout master && \
    cd src && \
    ./make.bash

ENV PATH="/go/goroot/bin:${PATH}"
