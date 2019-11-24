// +build js,wasm

package main

import (
	"net/url"
	"syscall/js"

	"github.com/cloudflare/cfssl/log"
)

func run(this js.Value, inputs []js.Value) interface{} {
	caCommonName := js.Global().Get("document").Call("getElementById", "root_common_name").Get("value").String()
	caCommonName = strOrDefault(caCommonName, "mtls.dev")

	caBundle, err := generateCACert(caCommonName)
	if err != nil {
		log.Errorf("generateCACert(%s): %s", caCommonName, err)
		return nil
	}

	serverCommonName := js.Global().Get("document").Call("getElementById", "server_common_name").Get("value").String()
	serverHostsJoined := js.Global().Get("document").Call("getElementById", "server_hosts").Get("value").String()
	serverExpiration := js.Global().Get("document").Call("getElementById", "server_duration").Get("value").String()
	serverConf := CertConfig{
		CommonName: strOrDefault(serverCommonName, "localhost"),
		Hosts:      strOrDefault(serverHostsJoined, "localhost,mtls.dev"),
		Expiration: strOrDefault(serverExpiration, defaultConfig.ExpiryString),
		CAPublic:   caBundle.Public,
		CAPrivate:  caBundle.Private,
	}

	serverBundle, err := generateServerCert(serverConf)
	if err != nil {
		log.Errorf("generateServerCert(%s, %s): %s", serverCommonName, serverHostsJoined, err)
		return nil
	}

	clientCommonName := js.Global().Get("document").Call("getElementById", "client_common_name").Get("value").String()
	clientExpiration := js.Global().Get("document").Call("getElementById", "client_duration").Get("value").String()
	clientConf := CertConfig{
		CommonName: strOrDefault(clientCommonName, "localhost"),
		Expiration: strOrDefault(clientExpiration, defaultConfig.ExpiryString),
		CAPublic:   caBundle.Public,
		CAPrivate:  caBundle.Private,
	}

	clientBundle, err := generateClientCert(clientConf)

	downloadAll(caBundle, serverBundle, clientBundle)

	return nil
}

func downloadAll(caBundle, serverBundle, clientBundle CertBundle) {
	download("ca-key.pem", caBundle.Private)
	download("ca.pem", caBundle.Public)
	download("server-key.pem", serverBundle.Private)
	download("server.pem", serverBundle.Public)
	download("client-key.pem", clientBundle.Private)
	download("client.pem", clientBundle.Public)
}

func download(filename string, contents []byte) {
	escaped := url.PathEscape(string(contents))
	document := js.Global().Get("document")
	a := document.Call("createElement", "a")
	a.Set("href", "data:text/plain;charset=utf-8,"+escaped)
	a.Set("download", filename)
	a.Set("style", "display: none;")
	document.Get("body").Call("appendChild", a)
	a.Call("click")
	document.Get("body").Call("removeChild", a)
}

func registerRunFunc() {
	js.Global().Set("run", js.FuncOf(run))
}

func strOrDefault(str string, def string) string {
	if str != "" {
		return str
	}

	return def
}
