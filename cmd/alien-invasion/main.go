package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dhubler/aliens"
)

var numAliens = flag.Int("numAliens", 10, "Number of aliens invading")
var numRounds = flag.Int("numRounds", 10000, "Limit the number of rounds the aliens perform before giving up")
var silent = flag.Bool("silent", false, "Supress log output but still output city report and fallen cities")
var seed = flag.Int64("seed", 0, "Optional random seed to control pseudo random results.  Default of zero for random each time")
var strict = flag.Bool("strict", false, "Use a more strict parse that does not back link any cities in opposite directions")
var outputFile = flag.String("outputFile", "", "Optional remaining cities output file")

func main() {
	var err error
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
		NumberAliens:   *numAliens,
		InvasionRounds: *numRounds,
		CityMapInput:   os.Stdin,
	}
	if *seed == 0 {
		options.Seed = time.Now().UnixNano()
	}
	if *outputFile != "" {
		out, err := os.Create(*outputFile)
		abortOnErr(err)
		defer func() {
			abortOnErr(out.Close())
		}()
		options.RemaingCitiesOutput = out
	} else {
		options.RemaingCitiesOutput = os.Stdout
	}

	err = aliens.Invade(options)
	abortOnErr(err)
}

func abortOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error runing invasion. %s\n", err.Error())
		os.Exit(1)
	}
}
