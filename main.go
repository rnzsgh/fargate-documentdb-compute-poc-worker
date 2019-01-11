package main

import (
	"flag"

	log "github.com/golang/glog"
)

func main() {

	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	log.Info("This is a test - we need compute")

	log.Flush()
}
