package aliens

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestSmallInvasion(t *testing.T) {
	// allow for pseudo-random invasions
	rnd = rand.New(rand.NewSource(0))

	mapRdr, err := os.Open("testdata/small-map.txt")
	if err != nil {
		t.Fatal(err)
	}
	cities, err := parse(mapRdr)
	if err != nil {
		t.Fatal(err)
	}
	aliens := createAliens(8)
	var buf bytes.Buffer
	invasion := &invasion{report: &buf}
	remaining := invasion.invade(cities, aliens, 0)
	err = dump(&buf, remaining)
	if err != nil {
		t.Fatal(err)
	}
	Golden(t, "testdata/small-map-invasion.golden", &buf)
}

func TestMediumInvasion(t *testing.T) {

	// allow for pseudo-random invasions
	rnd = rand.New(rand.NewSource(0))

	cities := generateCityMap(10)
	aliens := createAliens(100)
	var buf bytes.Buffer
	invasion := &invasion{report: &buf}
	remaining := invasion.invade(cities, aliens, 10001)
	err := dump(&buf, remaining)
	if err != nil {
		t.Fatal(err)
	}
	Golden(t, "testdata/medium-invasion.golden", &buf)
}

func generateCityMap(levels int) map[string]*city {
	root := &city{Name: "x"}
	pool := make(map[string]*city)
	pool[root.Name] = root
	generateCityMapNest(5, root, pool)
	return pool
}

func generateCityMapNest(levels int, parent *city, pool map[string]*city) {
	if levels == 0 {
		return
	}
	var added []*city
	for direction, label := range directionLabels {
		if parent.neighoringCity(direction) == nil {
			neighbor := &city{Name: fmt.Sprintf("%s%s", parent.Name, label[:1])}
			pool[neighbor.Name] = neighbor
			parent.addNeighbor(direction, neighbor)
			added = append(added, neighbor)
		}
	}
	for _, child := range added {
		generateCityMapNest(levels-1, child, pool)
	}
}
