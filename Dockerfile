FROM golang:1.13.4

RUN git clone https://go.googlesource.com/go goroot && \
    cd goroot && \
    git checkout master && \
    cd src && \
    ./all.bash

ENV PATH="/go/goroot/bin:${PATH}"
