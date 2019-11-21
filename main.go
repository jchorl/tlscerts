package main

import (
	"flag"
)

func main() {
	flag.Parse()

	registerRunFunc()

	done := make(chan bool)

	<-done
}
