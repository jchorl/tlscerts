// +build js,wasm

package main

import (
	"net/url"
	"syscall/js"

	"github.com/cloudflare/cfssl/log"
)

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
