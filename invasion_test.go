package aliens

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"testing"
)

var updateFlag = flag.Bool("update", false, "update expected golden file(s)")

func TestInvasions(t *testing.T) {
	tests := []struct {
		seed     int64
		cities   string
		expected string
	}{
		{
			seed:     1657964729860941318,
			cities:   "testdata/small-map.txt",
			expected: "testdata/aliens-lose.golden",
		},
		{
			seed:     1657982898578641344,
			cities:   "testdata/small-map.txt",
			expected: "testdata/aliens-oscilating.golden",
		},
		{
			seed:     10,
			cities:   "testdata/small-map.txt",
			expected: "testdata/aliens-trapped.golden",
		},
		{
			seed: 10,
			// if you go north enough, even south is eventually north on a sphere
			cities:   "testdata/circular-map.txt",
			expected: "testdata/nothing-left.golden",
		},
	}
	var buf bytes.Buffer
	resetLog := divertGlobalLogger(&buf)
	defer resetLog()
	for _, test := range tests {
		buf.Reset()
		in, err := os.Open(test.cities)
		if err != nil {
			t.Fatal(err)
		}
		err = Invade(Options{
			Seed:                test.seed,
			RemaingCitiesOutput: &buf,
			NumberAliens:        10,
			InvasionRounds:      10,
			CityMapInput:        in,
		})
		in.Close()
		Golden(t, *updateFlag, test.expected, &buf)
	}
}

// for unit test that want to capture and verify log output.  be sure to call returned
// function in defer
func divertGlobalLogger(capture io.Writer) func() {
	orig := log.Default().Flags()
	log.SetFlags(0)
	log.SetOutput(capture)
	return func() {
		log.SetFlags(orig)
		log.SetOutput(os.Stderr)
	}
}

func TestMediumInvasion(t *testing.T) {
	var buf bytes.Buffer
	resetLog := divertGlobalLogger(&buf)
	defer resetLog()
	invasion := &Invasion{
		rnd:    rand.New(rand.NewSource(0)),
		cities: generateCityMap(10),
		aliens: createAliens(100),
		rounds: 200,
	}
	invasion.invade()
	err := dump(&buf, invasion.remaining)
	if err != nil {
		t.Fatal(err)
	}
	Golden(t, *updateFlag, "testdata/medium-invasion.golden", &buf)
}

// generate a city map of a given level
//    North: ^     East: >
//    South: v     West: <
// Example
//                         .
//                 ------------------
//                 |    |     |     |
//                .^    .v    .>    .<
//                 |     |
//    -------------    ------...
//    |    |   |      |     |
//   .^^  .^> .^<    .vv  .>  ...
func generateCityMap(levels int) map[string]*city {
	root := &city{Name: "."}
	pool := make(map[string]*city)
	pool[root.Name] = root
	generateCityMapNest(5, root, pool)
	return pool
}

// recursive function to help generate city map
func generateCityMapNest(levels int, parent *city, pool map[string]*city) {
	if levels == 0 {
		return
	}
	var added []*city

	// pass 1 : add all the children first before recursing otherwise neighbors will
	// be different
	labels := []string{"^", "v", ">", "<"}
	for direction := range directions {
		if parent.neighoringCity(direction) == nil {
			neighbor := &city{Name: fmt.Sprintf("%s%s", parent.Name, labels[direction])}
			pool[neighbor.Name] = neighbor
			parent.addNeighbor(direction, neighbor)
			neighbor.addNeighbor(oppositeDirection(direction), parent)
			added = append(added, neighbor)
		}
	}
	// pass 2 : recurse on any new children
	for _, child := range added {
		generateCityMapNest(levels-1, child, pool)
	}
}
