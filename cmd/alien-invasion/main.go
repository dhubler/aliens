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
var numRounds = flag.Int("numRounds", 10, "Number of rounds the aliens perform before giving up")
var silent = flag.Bool("silent", false, "Supress log output but still output city report and fallen cities")

func main() {
	cl := flag.CommandLine
	cl.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: < city-map-file > report\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *silent {
		log.SetOutput(ioutil.Discard)
	}
	options := aliens.Options{
		ReportOutput:   os.Stdout,
		CityMapInput:   os.Stdin,
		NumberAliens:   *numAliens,
		InvasionRounds: *numRounds,
	}
	err := aliens.Invasion(options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error runing invasion. %s\n", err.Error())
		os.Exit(1)
	}
}
