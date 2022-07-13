package aliens

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// parse returns sorted array of cities by parsing input stream acoording to
// a specific format. see README.md for full spec
// Example:
//   Boston south=NewYork west=Albany
//   Albany east=Boston
//   ..
func parse(r io.Reader) (map[string]*city, error) {
	lines := bufio.NewReader(r)
	refs := make([]*cityRef, 0)
	for {
		line, err := lines.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		ref, err := parseCityRef(line)
		if err != nil {
			return nil, err
		}
		if ref != nil {
			refs = append(refs, ref)
		}
	}
	cities := make(map[string]*city, len(refs))
	// pass 1 : make cities
	for _, ref := range refs {
		if _, hasExisting := cities[ref.Name]; !hasExisting {
			cities[ref.Name] = &city{Name: ref.Name}
		}

		// assumption: do not require all neighbors to have dedicated line in
		// map. could reduce allocations by first checking if city
		// already exists but this bit simpler
		for direction := range directions {
			neighbor := ref.neighoringCity(direction)
			if neighbor != "" {
				if _, hasExisting := cities[neighbor]; !hasExisting {
					cities[neighbor] = &city{Name: neighbor}
				}
			}
		}
	}
	// pass 2 : build cities pointers in all directions
	var directionLabels = []string{"north", "south", "east", "west"}
	for _, ref := range refs {
		city := cities[ref.Name]
		for direction, directionLabel := range directionLabels {
			neighborName := ref.neighoringCity(direction)
			if neighborName == "" {
				continue
			}
			neighbor, found := cities[neighborName]
			if !found {
				return nil, fmt.Errorf("parse error. city '%s' neighbor '%s' to the '%s' was not found", ref.Name, neighborName, directionLabel)
			}
			city.addNeighbor(direction, neighbor)
		}
	}

	return cities, nil
}

func parseCityRef(line string) (*cityRef, error) {
	segs := strings.Split(strings.Trim(line, " \n"), " ")
	if len(segs) < 1 {
		return nil, fmt.Errorf("parse error, no city defined in '%s'", line)
	}
	name := segs[0]
	// opinion: ignore and allow blank lines
	if name == "" {
		return nil, nil
	}
	ref := cityRef{Name: name}
	for i := 1; i < len(segs); i++ {
		directionAndCity := strings.Split(segs[i], "=")
		if len(directionAndCity) != 2 {
			return nil, fmt.Errorf("parse error, invalid direction=city '%s'", segs[i])
		}
		// opinion: allows for redundant directions and takes last value
		switch directionAndCity[0] {
		case "north":
			ref.North = directionAndCity[1]
		case "south":
			ref.South = directionAndCity[1]
		case "east":
			ref.East = directionAndCity[1]
		case "west":
			ref.West = directionAndCity[1]
		default:
			return nil, fmt.Errorf("parse error, '%s' is not a recognized direction", directionAndCity[0])
		}
	}
	return &ref, nil
}

func dump(wtr io.Writer, cities map[string]*city) error {
	names := cityNames(cities)
	for _, name := range names {
		city := cities[name]
		if _, err := fmt.Fprint(wtr, name); err != nil {
			return err
		}
		for direction, directionLabel := range directionLabels {
			neighbor := city.neighoringCity(direction)
			if neighbor != nil {
				if _, err := fmt.Fprintf(wtr, " %s=%s", directionLabel, neighbor.Name); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprintln(wtr); err != nil {
			return err
		}
	}
	return nil
}
