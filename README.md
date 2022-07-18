# Alien Invasion Simulation

Simulate aliens invading a planet and report back the remaining cities.

Requirements:
* City maps are provided in files in a specific format detailed [here](#cityMapFormat)
* Any remaining cites at end of simulation should be reported in same format as input city maps.
* As cities fall to aliens, they are reported in specfic format detailed [here](#reportFallenCityFormat)
* Ability to control of the number of aliens in simulation

Simulation will honor the following aliens characteristics:

* If two aliens enter a city, the city MUST be destroyed along with the aliens.
* Once all aliens have landed in various cities, in the next round they MUST move to a neighboring city if there is one.  They MUST NOT stay in the same city. If there is no neighboring city to move to, they would be trapped in that city for the duration of the simulation.
* Destroyed cites MUST be removed from the map of available cities reported at the end of the simulation.
* On initial landing, if a city is already destroyed then incoming aliens MUST NOT land in that destroyed city.  This honors above requirement.
* If all remaining aliens are dead, simulation MAY end
* If no cities remain, simulation MAY end
* If all remaining aliens are trapped, simulation MAY end
* Aliens from a previous round in a given city MUST NOT interact with incoming aliens in the next round. For example, if in round `2`, alien `10` was in a city `X`.  Then in round `3` alien `99` enters city `X`.  This will not trigger a fight between aliens `10` and `99`.

Items of note:
* Because of the above requirements, aliens might oscilate between two cities until the end of the simulation. This happens when two cities only have one remaining exit path and that is to eachother. For example, if the only way out of Boston is to Bangor and the only way out of Bangor is to Boston and each city has a single alien then the aliens will constantly pass eachother in each round until the end of the simulation. 


# Setup

[Golang should be installed](https://go.dev/dl/) in PATH.  Any version should work but testing with Go v1.18 on Ubuntu 20.4

# Usage

```
cd cmd/alien-invasion
go run < ../../testdata/small-map.txt
```

Sample Output:

```
2022/07/13 20:04:19 invasion starting round
2022/07/13 20:04:19 alien 0 invading NewYork
2022/07/13 20:04:19 alien 1 invading NewYork
NewYork has been destroyed by alien 1 and alien 0!
2022/07/13 20:04:19 alien 4 invading Boston
2022/07/13 20:04:19 alien 7 invading Albany
2022/07/13 20:04:19 alien 8 invading Boston
Boston has been destroyed by alien 8 and alien 4!
2022/07/13 20:04:19 invasion 1 round
2022/07/13 20:04:19 no more aliens
Albany
Bangor
Columbus
Trenton
```

# Usage Options

```
Usage of ./alien-invasion: < city-map-file > report
  -numAliens int
    	Number of aliens invading (default 10)
  -numRounds int
    	Limit the number of rounds the aliens perform before giving up (default 10000)
  -outputFile string
    	Optional remaining cities output file
  -seed int
    	Optional random seed to control pseudo random results.  Default of zero for random each time
  -silent
    	Supress log output but still output city report and fallen cities
  -strict
    	Use a more strict parse that does not back link any cities in opposite directions
```

# Unit Testing

```
go test .
```

# <a name="cityMapFormat"></a>City map data format specification

Sample city input file:
```
Boston north=Bangor south=NewYork west=Albany
NewYork south=Trenton west=Columbus
```

Format Assumptions:

* City names cannot contain spaces
* If a city references another city, that referenced city **is not required** to have a separate line.  So in `Boston north=Bangor` then `Bangor south=Boston` is not required
* If a map contains inconsistent data with regard to neighboring references then those inconstencies will not be detected and results will be suspect.
Example of bad data:

```
Boston north=Bangor
Bangor south=Portland
```
# <a name="reportFallenCityFormat"></a>Fallen city format specification    

When a city falls to aliens, the city and the responsible aliens are reported in this format:

```
NewYork has been destroyed by alien 1 and alien 0!
```

Where `NewYork` is the city name. Aliens `1` and `0` are the responsible aliens.

# Developer Note - [Golden Files](https://ieftimov.com/posts/testing-in-go-golden-files/) in Unit Testing

Golden files are used to ensure large datasets only change when desired and in precise ways. If a unit test fails because the output doesn't match a "golden file" there are two options.  First inspect the "diff" and if the difference is expected, simply accept the difference by running the test again with the `-update` flag.  This strategy is used in the Golang SDK but not exclusive any single computer language.

Example:
```
go test -run TestMediumInvasion

        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -4,3 +4,2 @@
        	            	 Trenton has been destroyed by alien 7 and alien 4!
        	            	-Bogus
        	            	 Boston
        	Test:       	TestSmallInvasion
```

If difference **is expected**, then run following command to update the golden file

```
go test -run TestMediumInvasion -update
```

If difference **is not expected**, then you found a bug in your code.

# Developer Note - Pseudo Random Unit Testing

Some unit tests set the random number seed to get consistent random values, or "pseudo random" values.  Together with golden files, very complex code with randomness can have a predictable output to test against.  In order for this to work, certain code may have to avoid iterating Golang maps as that iteration is not deterministic.

# Developer Note - Test Coverage

To inspect unit test coverage:

```
go test -coverprofile cp.out .
go tool cover -html=cp.out
```