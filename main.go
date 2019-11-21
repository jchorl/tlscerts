package main

import (
	"flag"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/log"
)

// CertBundle just packages up a public cert and private key together
type CertBundle struct {
	Public  []byte
	Private []byte
}

func generateCACert(commonName string, hosts []string) (CertBundle, error) {
	req := csr.CertificateRequest{
		CN:         commonName,
		Hosts:      hosts,
		KeyRequest: csr.NewKeyRequest(),
	}

	cert, _, key, err := initca.New(&req)
	if err != nil {
		return CertBundle{}, fmt.Errorf("initca.New(...): %w", err)
	}

	return CertBundle{Public: cert, Private: key}, nil
}

func run(this js.Value, inputs []js.Value) interface{} {
	caCommonName := js.Global().Get("document").Call("getElementById", "root_common_name").Get("value").String()
	caHostsJoined := js.Global().Get("document").Call("getElementById", "root_hosts").Get("value").String()
	caHosts := strings.Split(caHostsJoined, ",")

	bundle, err := generateCACert(caCommonName, caHosts)
	if err != nil {
		log.Errorf("generateCACert(%s, %v): %s", caCommonName, caHosts, err)
	}

	log.Infof("public: %s", bundle.Public)
	log.Infof("private: %s", bundle.Private)

	return nil
}

func main() {
	flag.Parse()

	js.Global().Set("run", js.FuncOf(run))

	done := make(chan bool)

	<-done
}
