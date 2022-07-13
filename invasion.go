package aliens

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
)

// rnd will allow controlling the seed so invasions are pseudo-random
// and can be deterministic and potentially useful for unit testing
var rnd = rand.New(rand.NewSource(time.Now().Unix()))

type Options struct {
	CityMapInput   io.Reader
	NumberAliens   int
	InvationRounds int
	ReportOutput   io.Writer
}

// Invasion interface to run invasion simulation.
func Invasion(options Options) error {
	cities, err := parse(options.CityMapInput)
	if err != nil {
		return err
	}
	attack := &invasion{options.ReportOutput}
	aliens := createAliens(options.NumberAliens)
	remaingCities := attack.invade(cities, aliens, options.InvationRounds)
	return dump(options.ReportOutput, remaingCities)
}

// invasion simulates aliens navigating a map of cities according to a set of
// rules outlined in README.md.
// returns list of cities that are remaining after invasion

type invasion struct {
	report io.Writer
}

// acoording to spec.  does not include initial round
const maxRounds = 10000

func (attack *invasion) invade(cities map[string]*city, aliens []alien, rounds int) map[string]*city {
	destroyedCities := make(map[string]*city)
	invadedCities := make(map[*city]alien)
	startCityNames := cityNames(cities)

	// start aliens in random cities, cities can be destroyed in this phase
	log.Print("invasion starting round")
	for _, alien := range aliens {
		cityIndex := rnd.Intn(len(cities))
		city := cities[startCityNames[cityIndex]]
		if _, alreadyDestroyed := destroyedCities[city.Name]; alreadyDestroyed {
			// avoid cities that were already destroyed in this initial round
			continue
		}
		attack.invadeCity(alien, city, invadedCities, destroyedCities)
	}

	if rounds > maxRounds {
		log.Printf("warning, limited to %d exceeds maximum rounds or %d", rounds, maxRounds)
		rounds = maxRounds
	}

	// move aliens around until rounds are done
	for i := 0; i < rounds; i++ {
		log.Printf("invasion %d round", i+1)
		currentCities := invadedCities
		currentCitiesNames := invadedCityNames(currentCities)
		invadedCities := make(map[*city]alien)
		trappedCount := 0
		for _, origCityName := range currentCitiesNames {
			origCity := cities[origCityName]
			alien := currentCities[origCity]
			city := attack.nextRandomCity(rnd, origCity)
			if city == nil {
				trappedCount++
				continue
			}
			attack.invadeCity(alien, city, invadedCities, destroyedCities)
		}
		if len(invadedCities) == 0 {
			log.Printf("no more aliens")
			break
		}
		if trappedCount == len(invadedCities) {
			log.Printf("all remaining aliens are trapped")
			break
		}
	}

	// remaining = original list - destroyed
	remaining := make(map[string]*city)
	for name, city := range cities {
		if _, destroyed := destroyedCities[name]; !destroyed {
			remaining[name] = city
		}
	}

	return remaining
}

func (attack *invasion) invadeCity(incomingAlien alien, targetCity *city, invadedCities map[*city]alien, destroyedCities map[string]*city) {
	log.Printf("alien %s invading %s", incomingAlien, targetCity.Name)
	if invadedAlien, isInvaded := invadedCities[targetCity]; isInvaded {
		destroyedCities[targetCity.Name] = targetCity
		delete(invadedCities, targetCity) // leaves aliens inside
		fmt.Fprintf(attack.report, "%s has been destroyed by alien %s and alien %s!\n", targetCity.Name, incomingAlien, invadedAlien)
		targetCity.destroy(incomingAlien, invadedAlien)
	} else {
		invadedCities[targetCity] = incomingAlien
	}
}

// nextRandomCity picks a random neighboring city or return nil if
// there are no cities left
func (attack *invasion) nextRandomCity(rnd *rand.Rand, c *city) *city {
	startCityIndex := rnd.Intn(len(directions))
	for i := 0; i < len(directions); i++ {
		candidateIndex := (startCityIndex + i) % len(directions)
		candidate := c.neighoringCity(candidateIndex)
		if candidate != nil {
			return candidate
		}
	}
	return nil
}
