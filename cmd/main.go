package main

import (
	"flag"

	"github.com/boypt/crontasker"
)

func main() {
	config := ""
	flag.StringVar(&config, "c", "cronconf.txt", "set configuration `file`")
	flag.Parse()
	crontasker.CronDaemon(config)
}
