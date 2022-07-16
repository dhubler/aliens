package aliens

import (
	"bytes"
	"flag"
	"fmt"
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
			seed:     1657895446613653668,
			cities:   "testdata/small-map.txt",
			expected: "testdata/aliens-lose.golden",
		},
		{
			seed:     1657964729860941318,
			cities:   "testdata/small-map.txt",
			expected: "testdata/aliens-trapped-and-oscilating.golden",
		},
	}
	for _, test := range tests {
		var buf bytes.Buffer
		in, err := os.Open(test.cities)
		if err != nil {
			t.Fatal(err)
		}
		err = Invade(Options{
			Seed:                test.seed,
			RemaingCitiesOutput: &buf,
			NumberAliens:        10,
			InvasionRounds:      maxRounds,
			CityMapInput:        in,
		})
		in.Close()
		Golden(t, *updateFlag, test.expected, &buf)
	}
}

func TestMediumInvasion(t *testing.T) {
	var buf bytes.Buffer
	invasion := &Invasion{
		rnd:    rand.New(rand.NewSource(0)),
		cities: generateCityMap(10),
		aliens: createAliens(100),
		rounds: maxRounds,
	}
	invasion.invade()
	err := dump(&buf, invasion.remaining)
	if err != nil {
		t.Fatal(err)
	}
	Golden(t, *updateFlag, "testdata/medium-invasion.golden", &buf)
}

// generate a city map of a given level
// Example
//                         x
//                 ------------------
//                 |    |     |     |
//                xn    xs    xe    xw
//    -------------    ------...
//    |   |   |   |    |   |
//   xnn xns xne xnw  xsn xss  ...
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
			added = append(added, neighbor)
		}
	}
	// pass 2 : recurse on any new children
	for _, child := range added {
		generateCityMapNest(levels-1, child, pool)
	}
}
