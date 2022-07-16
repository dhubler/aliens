package aliens

import (
	"fmt"
	"io"
	"log"
	"math/rand"
)

// acoording to spec.  does not include initial round
const maxRounds = 10000

// Options to control the invasion
type Options struct {
	NumberAliens   int // how many aliens to start invasion
	InvasionRounds int // rounds to go before giving up

	// controls random seed so invasions are pseudo-random
	// and therefore deterministic and potentially useful for unit testing
	// or reproducing a particular invasion
	Seed int64

	CityMapInput        io.Reader
	RemaingCitiesOutput io.Writer
}

// NewInvasion interface to run invasion simulation.
func Invade(options Options) error {
	invasion := &Invasion{
		aliens: createAliens(options.NumberAliens),
		rnd:    rand.New(rand.NewSource(options.Seed)),
		rounds: options.InvasionRounds,
	}
	log.Printf("using random seed %d", options.Seed)
	fmt.Printf("using random seed %d", options.Seed)
	var err error
	invasion.cities, err = parse(options.CityMapInput)
	if err != nil {
		return err
	}
	invasion.invade()
	return dump(options.RemaingCitiesOutput, invasion.remaining)
}

type Invasion struct {
	rnd       *rand.Rand
	cities    map[string]*city
	remaining map[string]*city
	aliens    []alien
	rounds    int
}

// invade simulates aliens navigating a map of cities according to a set of
// rules outlined in README.md.
// returns list of cities that are remaining after invasion
func (attack *Invasion) invade() {
	destroyedCities := make(map[string]*city)
	invadedCities := make(map[*city]alien)
	trappedAlienCities := make(map[*city]alien)
	startCityNames := cityNames(attack.cities)

	if attack.rounds > maxRounds {
		log.Printf("warning, limited to %d exceeds maximum rounds or %d", attack.rounds, maxRounds)
		attack.rounds = maxRounds
	}

	// start aliens in random cities, cities can be destroyed in this phase
	log.Print("invasion starting round")
	for _, alien := range attack.aliens {
		cityIndex := attack.rnd.Intn(len(attack.cities))
		city := attack.cities[startCityNames[cityIndex]]
		if _, alreadyDestroyed := destroyedCities[city.Name]; alreadyDestroyed {
			// avoid landing in cities that were already destroyed in this initial round
			continue
		}
		attack.invadeCity(alien, city, invadedCities, destroyedCities)
	}

	// move aliens around until rounds are done
	for i := 0; i < attack.rounds; i++ {
		log.Printf("invasion %d round", i+1)
		currentCities := invadedCities

		// we iterate the sorted city names to allow for pseudom random test
		// cases.  Otherwise iterating invadedCities would be bit faster and
		// simpler
		currentCitiesNames := invadedCityNames(currentCities)

		invadedCities = make(map[*city]alien)
		for _, origCityName := range currentCitiesNames {
			origCity := attack.cities[origCityName]
			alien := currentCities[origCity]
			city := attack.nextRandomCity(origCity)
			if city == nil {
				trappedAlienCities[origCity] = alien
			} else {
				attack.invadeCity(alien, city, invadedCities, destroyedCities)
			}
		}
		if len(invadedCities) == 0 {
			break
		}
	}

	log.Printf("%d alien(s) left, %d alien(s) trapped", len(invadedCities), len(trappedAlienCities))

	// remaining = original list - destroyed
	attack.remaining = make(map[string]*city)
	for name, city := range attack.cities {
		if _, destroyed := destroyedCities[name]; !destroyed {
			attack.remaining[name] = city
		}
	}
}

// invadeCity checks if another alien is in city to trigger a destroy or if this
// is just first visit
func (attack *Invasion) invadeCity(incomingAlien alien, targetCity *city, invadedCities map[*city]alien, destroyedCities map[string]*city) {
	log.Printf("alien %s invading %s", incomingAlien, targetCity.Name)
	if invadedAlien, isInvaded := invadedCities[targetCity]; isInvaded {
		destroyedCities[targetCity.Name] = targetCity
		delete(invadedCities, targetCity) // leaves aliens inside
		log.Printf("%s has been destroyed by alien %s and alien %s!\n", targetCity.Name, incomingAlien, invadedAlien)
		targetCity.destroy(incomingAlien, invadedAlien)
	} else {
		invadedCities[targetCity] = incomingAlien
	}
}

// nextRandomCity picks a random neighboring city or return nil if
// there are no cities left
func (attack *Invasion) nextRandomCity(c *city) *city {
	startCityIndex := attack.rnd.Intn(len(directions))
	for i := 0; i < len(directions); i++ {
		candidateIndex := (startCityIndex + i) % len(directions)
		candidate := c.neighoringCity(candidateIndex)
		if candidate != nil {
			return candidate
		}
	}
	return nil
}
