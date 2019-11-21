package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	cert, err := tls.LoadX509KeyPair("./client.pem", "./client-key.pem")
	if err != nil {
		log.Fatalf("loading key pair: %s", err)
	}

	caCert, err := ioutil.ReadFile("./ca.pem")
	if err != nil {
		log.Fatalf("reading ca cert: %s", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: conf,
		},
	}

	resp, err := client.Get("https://server/")
	if err != nil {
		log.Fatalf("sending request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}
}
