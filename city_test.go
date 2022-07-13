package aliens

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCityRef(t *testing.T) {
	r := cityRef{Name: "x", North: "n", South: "s", East: "e", West: "w"}
	assert.Equal(t, "n", r.neighoringCity(North))
	assert.Equal(t, "s", r.neighoringCity(South))
	assert.Equal(t, "e", r.neighoringCity(East))
	assert.Equal(t, "w", r.neighoringCity(West))
}

func TestCity(t *testing.T) {
	x := &city{Name: "x"}
	n := &city{Name: "n"}
	x.addNeighbor(North, n)
	s := &city{Name: "s"}
	x.addNeighbor(South, s)
	e := &city{Name: "e"}
	x.addNeighbor(East, e)
	w := &city{Name: "w"}
	x.addNeighbor(West, w)

	t.Run("neighbor", func(t *testing.T) {
		assert.Equal(t, s, x.neighoringCity(South))
		assert.Equal(t, n, x.neighoringCity(North))
		assert.Equal(t, e, x.neighoringCity(East))
		assert.Equal(t, w, x.neighoringCity(West))
	})

	t.Run("destroy", func(t *testing.T) {
		x.destroy(alien("a"), alien("b"))
		assert.Nil(t, x.South)
		assert.Nil(t, s.North)
	
		assert.Nil(t, x.North)
		assert.Nil(t, n.South)
	
		assert.Nil(t, x.West)
		assert.Nil(t, w.East)
	
		assert.Nil(t, x.East)
		assert.Nil(t, e.West)	
	})
}
