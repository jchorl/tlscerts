// +build !js,!wasm

package main

import "github.com/cloudflare/cfssl/log"

func registerRunFunc() {
	log.Info("should register run function, but running outside of js env")
}
