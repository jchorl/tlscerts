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
		CAPublic:   caBundle.Public,
		CAPrivate:  caBundle.Private,
	})
	require.NoError(t, err)

	clientBundle, err := generateServerCert(CertConfig{
		CommonName: "mtls.dev",
		CAPublic:   caBundle.Public,
		CAPrivate:  caBundle.Private,
	})
	require.NoError(t, err)

	serverTLSConf := getTLSConfig(t, serverBundle.Public, serverBundle.Private, caBundle.Public, true)

	srv := newTestServer(t, serverTLSConf)
	defer srv.Close()

	clientTLSConf := getTLSConfig(t, clientBundle.Public, clientBundle.Private, caBundle.Public, false)

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

func getTLSConfig(t *testing.T, public, private, caPublic []byte, isServer bool) *tls.Config {
	cert, err := tls.X509KeyPair(public, private)
	require.NoError(t, err)

	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caPublic)

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
