package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dhubler/aliens"
)

var numAliens = flag.Int("numAliens", 10, "Number of aliens invading")
var numRounds = flag.Int("numRounds", 10, "Number of rounds the aliens perform beforing giving up")
var silent = flag.Bool("silent", false, "No log output but still output city report and fallen cities")

func main() {
	flag.Parse()
	if *silent {
		log.SetOutput(ioutil.Discard)
	}
	options := aliens.Options{
		ReportOutput:   os.Stdout,
		CityMapInput:   os.Stdin,
		NumberAliens:   *numAliens,
		InvationRounds: *numRounds,
	}
	err := aliens.Invasion(options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error runing invasion. %s\n", err.Error())
		os.Exit(1)
	}
}
