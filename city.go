package aliens

import (
	"fmt"
	"sort"
)

// neighboring directions supported in each city
const (
	North int = iota
	South
	East
	West
)

// convient list of directions in order of constants.
var directions = []int{
	North, South, East, West,
}

// used when encoding and decoding maps
var directionLabels = []string{
	"north", "south", "east", "west",
}

type city struct {
	Name  string
	North *city
	South *city
	East  *city
	West  *city
}

// cityNames are in city name sorted order
func cityNames(cities map[string]*city) []string {
	names := make([]string, 0, len(cities))
	for name := range cities {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// invadedCityNames are in city name sorted order
func invadedCityNames(cities map[*city]alien) []string {
	names := make([]string, 0, len(cities))
	for city := range cities {
		names = append(names, city.Name)
	}
	sort.Strings(names)
	return names
}

// addNeighbor will add a neighboring city to a given city AND
// will also add given city as a reference back to the given city
// in the opposite compass direction. e.g. If you add a neighbor
// to the south, that neighbor will have a neighbor to the north
// to the original city
func (c *city) addNeighbor(direction int, neighbor *city) {
	switch direction {
	case North:
		c.North = neighbor
		neighbor.South = c
	case South:
		c.South = neighbor
		neighbor.North = c
	case West:
		c.West = neighbor
		neighbor.East = c
	case East:
		c.East = neighbor
		neighbor.West = c
	default:
		panic(fmt.Errorf("invalid direction %d", direction))
	}
}

// neighoringCity gets a neighbor in a specific direction.  If the city doesn't
// have a neighbor in that direction, nil is returned
func (c *city) neighoringCity(direction int) *city {
	switch direction {
	case North:
		return c.North
	case South:
		return c.South
	case West:
		return c.West
	case East:
		return c.East
	default:
		panic(fmt.Errorf("invalid direction %d", direction))
	}
}

// destroy will trap any aliens in the city and remove all roads
// into city from neighboring cities
func (c *city) destroy(a, b alien) {
	for direction := range directions {
		neighbor := c.neighoringCity(direction)
		if neighbor != nil {
			// remove all pointers back to destroyed city from neighboring cities
			// where pointer is defined by reference in the opposite direction
			switch direction {
			case North:
				c.North = nil
				neighbor.South = nil
			case South:
				c.South = nil
				neighbor.North = nil
			case East:
				c.East = nil
				neighbor.West = nil
			case West:
				c.West = nil
				neighbor.East = nil
			}
		}
	}
}

// cityRef is a temporary struct used as a holding place to ultimately
// build city map
type cityRef struct {
	Name  string
	North string
	South string
	East  string
	West  string
}

func (c cityRef) neighoringCity(direction int) string {
	switch direction {
	case North:
		return c.North
	case South:
		return c.South
	case West:
		return c.West
	case East:
		return c.East
	default:
		panic(fmt.Errorf("invalid direction %d", direction))
	}
}
