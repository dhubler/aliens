package aliens

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicParse(t *testing.T) {
	in := `
Bar south=Foo west=Bee
Foo north=Bar south=Qu-ux west=Baz
`
	cityMap, err := parse(strings.NewReader(in))
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

func TestLargeDump(t *testing.T) {
	pool := generateCityMap(5)
	var actual bytes.Buffer
	err := dump(&actual, pool)
	assert.NoError(t, err)
	Golden(t, *updateFlag, "testdata/large-dump-map.golden", &actual)
}
