package main

import (
	"flag"
	"time"

	log "github.com/golang/glog"
)

func main() {

	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	log.Info("This is a test - we need compute")

	time.Sleep(3 * time.Minute)

	log.Flush()
}
