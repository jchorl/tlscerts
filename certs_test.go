package main

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCertExpiryServer(t *testing.T) {
	caBundle, err := generateCACert("CA")
	require.NoError(t, err)

	serverBundle, err := generateServerCert(CertConfig{
		CommonName: "mtls.dev",
		Hosts:      "127.0.0.1",
		Expiration: "127h",
		CACert:     caBundle.Cert,
		CAKey:      caBundle.Key,
	})
	require.NoError(t, err)

	certDERBlock, _ := pem.Decode(serverBundle.Cert)
	cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	require.NoError(t, err)

	// check that the expiry is ~127 hours from now
	require.True(t, cert.NotAfter.Sub(time.Now()) > 126*time.Hour)
	require.True(t, cert.NotAfter.Sub(time.Now()) < 128*time.Hour)
}

func TestCertExpiryClient(t *testing.T) {
	caBundle, err := generateCACert("CA")
	require.NoError(t, err)

	clientBundle, err := generateClientCert(CertConfig{
		CommonName: "mtls.dev",
		Hosts:      "127.0.0.1",
		Expiration: "127h",
		CACert:     caBundle.Cert,
		CAKey:      caBundle.Key,
	})
	require.NoError(t, err)

	certDERBlock, _ := pem.Decode(clientBundle.Cert)
	cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	require.NoError(t, err)

	// check that the expiry is ~127 hours from now
	require.True(t, cert.NotAfter.Sub(time.Now()) > 126*time.Hour)
	require.True(t, cert.NotAfter.Sub(time.Now()) < 128*time.Hour)
}
