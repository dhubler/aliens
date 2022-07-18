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
	assert.NoError(t, x.addNeighbor(North, n))
	s := &city{Name: "s"}
	assert.NoError(t, x.addNeighbor(South, s))
	e := &city{Name: "e"}
	assert.NoError(t, x.addNeighbor(East, e))
	w := &city{Name: "w"}
	assert.NoError(t, x.addNeighbor(West, w))

	t.Run("harmless neighbors", func(t *testing.T) {
		assert.NoError(t, x.addNeighbor(South, s))
	})

	t.Run("bad neighbor  ", func(t *testing.T) {
		assert.Error(t, x.addNeighbor(North, s))
	})

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

func TestAddNeighborBidirectional(t *testing.T) {
	a := &city{Name: "a"}
	b := &city{Name: "b"}
	c := &city{Name: "c"}
	err := a.addNeighborBidiectional(North, b)
	assert.NoError(t, err)
	err = b.addNeighborBidiectional(South, c)
	assert.Error(t, err)
}
