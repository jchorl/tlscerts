package main

import (
	"fmt"
	"time"

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
	Cert []byte
	Key  []byte
}

// CertConfig lays out some config options for generating a cert
type CertConfig struct {
	CommonName string
	Hosts      string
	Expiration string
	CACert     []byte
	CAKey      []byte
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

	return CertBundle{Cert: cert, Key: key}, nil
}

func generateServerCert(cfg CertConfig) (CertBundle, error) {
	hosts := signer.SplitHosts(cfg.Hosts)
	srvConfig := DefaultServerConfig()
	if cfg.Expiration != "" {
		parsed, err := time.ParseDuration(cfg.Expiration)
		if err != nil {
			return CertBundle{}, fmt.Errorf("invalid duration %s: %w", cfg.Expiration, err)
		}
		srvConfig.Signing.Profiles["www"].Expiry = parsed
		srvConfig.Signing.Profiles["www"].ExpiryString = parsed.String()
	}

	return generateCert(cfg.CommonName, hosts, "www", srvConfig.Signing, cfg.CACert, cfg.CAKey)
}

func generateClientCert(cfg CertConfig) (CertBundle, error) {
	cliConfig := DefaultClientConfig()
	if cfg.Expiration != "" {
		parsed, err := time.ParseDuration(cfg.Expiration)
		if err != nil {
			return CertBundle{}, fmt.Errorf("invalid duration %s: %w", cfg.Expiration, err)
		}
		cliConfig.Signing.Profiles["client"].Expiry = parsed
		cliConfig.Signing.Profiles["client"].ExpiryString = parsed.String()
	}
	return generateCert(cfg.CommonName, nil, "client", cliConfig.Signing, cfg.CACert, cfg.CAKey)
}

func generateCert(commonName string, hosts []string, profile string, signingConfig *config.Signing, ca []byte, caKey []byte) (CertBundle, error) {
	log.Infof("received cert generate request, commonName=%s, hosts=%v, profile=%s, signingConfig=%+v", commonName, hosts, profile, signingConfig)
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
		Cert: cert,
		Key:  key,
	}, nil
}
