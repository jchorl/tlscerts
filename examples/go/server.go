package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	caCert, err := ioutil.ReadFile("./ca.pem")
	if err != nil {
		log.Fatalf("reading ca cert: %s", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	conf := &tls.Config{
		ClientCAs: certPool,
	}

	conf.BuildNameToCertificate()

	httpServer := &http.Server{
		Addr:      ":443",
		TLSConfig: conf,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	err = httpServer.ListenAndServeTLS("./server.pem", "./server-key.pem")
	if err != nil {
		log.Fatalf("listen and serving: %s", err)
	}
}
