# Alien Invasion Simulation

Simulate aliens invading a planet and report back the remaining cities.

Requirements:
* City maps are provided in files in a specific format detailed [here](#cityMapFormat)
* Any remaining cites at end of simulation should be reported in same format as input city maps.
* As cities fall to aliens, they are reported in specfic format detailed [here](#reportFallenCityFormat)
* Ability to control of the number of aliens in simulation

Simulation will honor the following aliens characteristics:

* If two aliens enter a city, city is destroyed along with the aliens.
* Once all aliens have landed in various cities, in the next round they move to a neighboring city if there is one.  If there isn't one, they would be trapped in that city for the duration of the simulation.
* Destroyed cites are no longer on the map of available cities
* On initial landing, if a city is already destroyed then incoming aliens will not land in that destroyed city keeping to the requirement that once a city is destroyed it is off the map.
* If all remaining aliens are dead, simulation will end
* If all remaining aliens are trapped, simulation will end
* Aliens all attempt to leave respective cities then enter a new random city one at a time. So if in round `2`, alien `10` was in a city `X`.  Then in round `3` alien `99` enters city `X`.  This will not trigger a fight between aliens `10` and `99`.

# Setup

[Golang should be installed](https://go.dev/dl/) in PATH.  Any version should work but testing with Go v1.18 on Ubuntu 20.4

# Usage

**Step 1.)** Compile binary

    cd cmd/alien-invasion
    go build

**Step 2.)** Run an invasion

    ./alien-invasion < ../../testdata/small-map.txt

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
 ./alien-invasion --help
 Usage of ./alien-invasion:
  -numAliens int
    	Number of aliens invading (default 10)
  -numRounds int
    	Number of rounds the aliens perform before giving up (default 10)
  -silent
    	Supress log output but still output city report and fallen cities
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