package main

import (
	"flag"
	"fmt"

	"github.com/boypt/crontasker"
)

func main() {
	config := ""
	flag.StringVar(&config, "c", "cronconf.txt", "set configuration `file`")
	flag.Parse()
	err := crontasker.CronDaemon(config)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
