package aliens

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicParse(t *testing.T) {
	in := `
Bar south=Foo west=Bee
Foo north=Bar south=Qu-ux west=Baz
`
	cityMap, err := parse(strings.NewReader(in), false)
	assert.NoError(t, err)
	assert.NotNil(t, cityMap)
	assert.Equal(t, 5, len(cityMap))
	var out bytes.Buffer
	err = dump(&out, cityMap)
	assert.NoError(t, err)
	expected := `Bar south=Foo west=Bee
Baz east=Foo
Bee east=Bar
Foo north=Bar south=Qu-ux west=Baz
Qu-ux north=Foo
`
	assert.Equal(t, expected, out.String())
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		line     string
		expected *cityRef
		invalid  bool
	}{
		{
			line:     "Foo",
			expected: &cityRef{Name: "Foo"},
		},
		{
			line:    "Foo north =Bar",
			invalid: true,
		},
		{
			line:    "Foo north=Bar bogus",
			invalid: true,
		},
		{
			line:    "Foo Bar north=Bar",
			invalid: true,
		},
		{
			line:    "Foo norf=Goo",
			invalid: true,
		},
		{
			line:    "Foo north==Goo",
			invalid: true,
		},
		{
			line:    "Foo north= south=Goo",
			invalid: true,
		},
		{
			line:     "Foo north=Bar",
			expected: &cityRef{Name: "Foo", North: "Bar"},
		},
		{
			line:     "Foo south=Bar",
			expected: &cityRef{Name: "Foo", South: "Bar"},
		},
		{
			line:     "Foo east=Bar",
			expected: &cityRef{Name: "Foo", East: "Bar"},
		},
		{
			line:     "Foo west=Bar",
			expected: &cityRef{Name: "Foo", West: "Bar"},
		},
		{
			line:     "Foo north=a south=b east=c west=d",
			expected: &cityRef{Name: "Foo", North: "a", South: "b", East: "c", West: "d"},
		},
	}
	for _, test := range tests {
		actual, err := parseCityRef(test.line)
		if test.expected != nil {
			assert.Equal(t, test.expected, actual, test.line)
		} else if test.invalid {
			assert.Error(t, err, test.line)
		}
	}
}

func TestParseMaps(t *testing.T) {
	tests := []struct {
		src            string
		expected       string
		expectedStrict string
	}{
		{
			"testdata/small-map.txt",
			"testdata/small-map.golden",
			"testdata/small-map-strict.golden",
		},
		{
			"testdata/circular-map.txt",
			"testdata/circular-map.golden",
			"testdata/circular-map-strict.golden",
		},
		{
			"testdata/bad-map.txt",
			"testdata/bad-map.golden",
			"testdata/bad-map-strict.golden",
		},
	}

	for _, test := range tests {
		src, err := ioutil.ReadFile(test.src)
		if err != nil {
			t.Fatal(err, test.src)
		}
		actual, err := parse(bytes.NewBuffer(src), false)
		assert.NoError(t, err)
		var buf bytes.Buffer
		dump(&buf, actual)
		Golden(t, *updateFlag, test.expected, &buf)

		actual2, err := parse(bytes.NewBuffer(src), true)
		assert.NoError(t, err)
		var buf2 bytes.Buffer
		dump(&buf2, actual2)
		Golden(t, *updateFlag, test.expectedStrict, &buf2)
	}
}

func TestLargeDump(t *testing.T) { // grow up
	pool := generateCityMap(5)
	var actual bytes.Buffer
	err := dump(&actual, pool)
	assert.NoError(t, err)
	Golden(t, *updateFlag, "testdata/large-dump-map.golden", &actual)
}
