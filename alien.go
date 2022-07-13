package aliens

import (
	"strconv"
)

type alien string

func createAliens(numAliens int) []alien {
	aliens := make([]alien, numAliens)
	for i := 0; i < numAliens; i++ {
		aliens[i] = alien(strconv.Itoa(i))
	}
	return aliens
}
