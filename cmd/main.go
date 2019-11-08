package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/boypt/crontasker"
)

var (
	config string
	debug  bool
	VER    string = "0.0.0-src"
)

func main() {
	flag.StringVar(&config, "c", "cronconf.txt", "set configuration `file`")
	flag.BoolVar(&debug, "debug", false, "show debug log")
	flag.Parse()
	fmt.Printf("crontasker ver: %s \n", VER)
	crontasker.SetDebug(debug)
	err := crontasker.CronDaemon(config)
	if err != nil {
		log.Fatal(err)
	}
}
