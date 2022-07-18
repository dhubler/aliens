package aliens

import (
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
	StrictMapParse      bool
}

// NewInvasion interface to run invasion simulation.
func Invade(options Options) error {
	invasion := &Invasion{
		aliens: createAliens(options.NumberAliens),
		rnd:    rand.New(rand.NewSource(options.Seed)),
		rounds: options.InvasionRounds,
	}
	log.Printf("using random seed %d", options.Seed)
	var err error
	invasion.cities, err = parse(options.CityMapInput, options.StrictMapParse)
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
func (sim *Invasion) invade() {
	destroyedCities := make(map[string]*city)
	invadedCities := make(map[*city]alien)
	trappedAlienCities := make(map[*city]alien)
	startCityNames := cityNames(sim.cities)

	if sim.rounds > maxRounds {
		log.Printf("warning, limited to %d exceeds maximum rounds or %d", sim.rounds, maxRounds)
		sim.rounds = maxRounds
	}

	// start aliens in random cities, cities can be destroyed in this phase
	log.Print("invasion starting round")
	for _, alien := range sim.aliens {
		cityIndex := sim.rnd.Intn(len(sim.cities))
	reattemptLanding:
		city := sim.cities[startCityNames[cityIndex]]
		if _, alreadyDestroyed := destroyedCities[city.Name]; alreadyDestroyed {
			// avoid landing in cities that were already destroyed in this initial round
			if len(destroyedCities) == len(sim.cities) {
				// no more cities to attack
				break
			}
			// go to next city, do not pick another random city because if there is
			// only 1 city left in a large list, finding it randomly would be inefficient
			cityIndex = (cityIndex + 1) % len(sim.cities)
			goto reattemptLanding
		}
		sim.invadeCity(alien, city, invadedCities, destroyedCities)
	}

	// move aliens around until rounds are done
	for i := 0; i < sim.rounds; i++ {
		log.Printf("invasion %d round", i+1)
		currentCities := invadedCities

		// we iterate the sorted city names to allow for pseudom random test
		// cases.  Otherwise iterating invadedCities would be bit faster and
		// simpler
		currentCitiesNames := invadedCityNames(currentCities)

		invadedCities = make(map[*city]alien)
		for _, origCityName := range currentCitiesNames {
			origCity := sim.cities[origCityName]
			alien := currentCities[origCity]
			city := sim.nextRandomCity(origCity)
			if city == nil {
				trappedAlienCities[origCity] = alien
			} else {
				sim.invadeCity(alien, city, invadedCities, destroyedCities)
			}
		}
		if len(invadedCities) == 0 {
			break
		}
	}

	// remaining = original list - destroyed
	sim.remaining = make(map[string]*city)
	for name, city := range sim.cities {
		if _, destroyed := destroyedCities[name]; !destroyed {
			sim.remaining[name] = city
		}
	}

	log.Printf("%d cities left, %d alien(s) left, %d alien(s) trapped", len(sim.remaining), len(invadedCities), len(trappedAlienCities))
}

// invadeCity checks if another alien is in city to trigger a destroy or if this
// is just first visit
func (sim *Invasion) invadeCity(incomingAlien alien, targetCity *city, invadedCities map[*city]alien, destroyedCities map[string]*city) {
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
func (sim *Invasion) nextRandomCity(c *city) *city {
	startCityIndex := sim.rnd.Intn(len(directions))
	for i := 0; i < len(directions); i++ {
		candidateIndex := (startCityIndex + i) % len(directions)
		candidate := c.neighoringCity(candidateIndex)
		if candidate != nil {
			return candidate
		}
	}
	return nil
}
