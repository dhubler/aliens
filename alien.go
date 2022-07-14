package aliens

import (
	"strconv"
)

// alien has no data except a name string
type alien string

// createAliens generates a set of aliens ready for invasion
func createAliens(numAliens int) []alien {
	aliens := make([]alien, numAliens)
	for i := 0; i < numAliens; i++ {
		// aliens have very boring names of a simple sequential
		// number but anything unique is allowed
		aliens[i] = alien(strconv.Itoa(i))
	}
	return aliens
}
