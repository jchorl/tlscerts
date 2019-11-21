// +build js,wasm

package main

import (
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

	return nil
}

func registerRunFunc() {
	js.Global().Set("run", js.FuncOf(run))
}
