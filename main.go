package main

import (
	"flag"
	"fmt"
	"syscall/js"

	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/log"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
)

// CertBundle just packages up a public cert and private key together
type CertBundle struct {
	Public  []byte
	Private []byte
}

func generateCACert(commonName string) (CertBundle, error) {
	req := csr.CertificateRequest{
		CN:         commonName,
		KeyRequest: csr.NewKeyRequest(),
	}

	cert, _, key, err := initca.New(&req)
	if err != nil {
		return CertBundle{}, fmt.Errorf("initca.New(...): %w", err)
	}

	return CertBundle{Public: cert, Private: key}, nil
}

func generateServerCert(commonName string, hostsJoined string, ca []byte, caKey []byte) (CertBundle, error) {
	hosts := signer.SplitHosts(hostsJoined)

	return generateCert(commonName, hosts, "www", ServerConfig.Signing, ca, caKey)
}

func generateClientCert(commonName string, ca []byte, caKey []byte) (CertBundle, error) {
	return generateCert(commonName, nil, "client", ClientConfig.Signing, ca, caKey)
}

func generateCert(commonName string, hosts []string, profile string, signingConfig *config.Signing, ca []byte, caKey []byte) (CertBundle, error) {
	req := csr.CertificateRequest{
		CN:         commonName,
		Hosts:      hosts,
		KeyRequest: csr.NewKeyRequest(),
	}

	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err := g.ProcessRequest(&req)
	if err != nil {
		key = nil
		return CertBundle{}, fmt.Errorf("g.ProcessRequest(%+v): %w", req, err)
	}

	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   hosts,
		Profile: profile,
	}

	parsedCa, err := helpers.ParseCertificatePEM(ca)
	if err != nil {
		return CertBundle{}, fmt.Errorf("helpers.ParseCertificatePEM(...): %w", err)
	}

	priv, err := helpers.ParsePrivateKeyPEMWithPassword(caKey, []byte{})
	if err != nil {
		return CertBundle{}, fmt.Errorf("helpers.ParsePrivateKeyPEMWithPassword(...): %w", err)
	}

	s, err := local.NewSigner(priv, parsedCa, signer.DefaultSigAlgo(priv), signingConfig)
	if err != nil {
		return CertBundle{}, fmt.Errorf("local.NewSigner(...): %w", err)
	}

	cert, err := s.Sign(signReq)
	if err != nil {
		return CertBundle{}, fmt.Errorf("s.Sign(...): %w", err)
	}

	return CertBundle{
		Public:  cert,
		Private: key,
	}, nil
}

func run(this js.Value, inputs []js.Value) interface{} {
	caCommonName := js.Global().Get("document").Call("getElementById", "root_common_name").Get("value").String()

	caBundle, err := generateCACert(caCommonName)
	if err != nil {
		log.Errorf("generateCACert(%s): %s", caCommonName, err)
		return nil
	}

	log.Infof("root public: %s", caBundle.Public)
	log.Infof("root private: %s", caBundle.Private)

	serverCommonName := js.Global().Get("document").Call("getElementById", "server_common_name").Get("value").String()
	serverHostsJoined := js.Global().Get("document").Call("getElementById", "server_hosts").Get("value").String()

	serverBundle, err := generateServerCert(serverCommonName, serverHostsJoined, caBundle.Public, caBundle.Private)
	if err != nil {
		log.Errorf("generateServerCert(%s, %s): %s", serverCommonName, serverHostsJoined, err)
		return nil
	}

	log.Infof("server public: %s", serverBundle.Public)
	log.Infof("server private: %s", serverBundle.Private)

	clientCommonName := js.Global().Get("document").Call("getElementById", "client_common_name").Get("value").String()
	clientBundle, err := generateClientCert(clientCommonName, caBundle.Public, caBundle.Private)

	log.Infof("client public: %s", clientBundle.Public)
	log.Infof("client private: %s", clientBundle.Private)

	return nil
}

func main() {
	flag.Parse()

	js.Global().Set("run", js.FuncOf(run))

	done := make(chan bool)

	<-done
}
