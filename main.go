package main

import (
	"flag"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/log"
)

func main() {
	flag.Parse()

	req := csr.CertificateRequest{
		CN:         "CA",
		Hosts:      []string{"tlscerts.dev"},
		KeyRequest: csr.NewKeyRequest(),
	}

	cert, csrPEM, key, err := initca.New(&req)
	if err != nil {
		log.Fatalf("initca.New(...): %s", err)
	}

	log.Infof("cert: %s", cert)
	log.Infof("csr: %s", csrPEM)
	log.Infof("key: %s", key)
}
