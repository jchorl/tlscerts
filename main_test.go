package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMTLS(t *testing.T) {
	caBundle, err := generateCACert("CA")
	require.NoError(t, err)

	serverBundle, err := generateServerCert(CertConfig{
		CommonName: "mtls.dev",
		Hosts:      "127.0.0.1",
		CACert:     caBundle.Cert,
		CAKey:      caBundle.Key,
	})
	require.NoError(t, err)

	clientBundle, err := generateServerCert(CertConfig{
		CommonName: "mtls.dev",
		CACert:     caBundle.Cert,
		CAKey:      caBundle.Key,
	})
	require.NoError(t, err)

	serverTLSConf := getTLSConfig(t, serverBundle.Cert, serverBundle.Key, caBundle.Cert, true)

	srv := newTestServer(t, serverTLSConf)
	defer srv.Close()

	clientTLSConf := getTLSConfig(t, clientBundle.Cert, clientBundle.Key, caBundle.Cert, false)

	client := srv.Client()
	client.Transport = &http.Transport{
		TLSClientConfig: clientTLSConf,
	}

	resp, err := client.Get(srv.URL)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, []byte("hi!\n"), body)
}

func getTLSConfig(t *testing.T, cert, key, caCert []byte, isServer bool) *tls.Config {
	pair, err := tls.X509KeyPair(cert, key)
	require.NoError(t, err)

	conf := &tls.Config{
		Certificates: []tls.Certificate{pair},
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	if isServer {
		conf.ClientCAs = certPool
	} else {
		conf.RootCAs = certPool
	}

	conf.BuildNameToCertificate()

	return conf
}

func newTestServer(t *testing.T, tlsConf *tls.Config) *httptest.Server {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hi!")
	}))
	ts.TLS = tlsConf
	ts.StartTLS()
	return ts
}
