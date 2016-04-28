package rure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	re, err := Compile(`(`)
	assert.Nil(t, re)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Unclosed parenthesis")
}

func TestIsMatch(t *testing.T) {
	re := MustCompile(`\p{So}`)
	assert.True(t, re.IsMatch("snowman: ☃"))
}

func TestFind(t *testing.T) {
	re := MustCompile(`\p{So}`)
	start, end, ok := re.Find("snowman: ☃")
	assert.True(t, ok)
	assert.Equal(t, 9, start)
	assert.Equal(t, 12, end)
}

func TestCaptures(t *testing.T) {
	re := MustCompile(`.(.*(?P<snowman>\p{So}))$`)
	caps := re.Captures("snowman: ☃")
	assert.NotNil(t, caps)

	start, end, ok := caps.Pos(2)
	assert.True(t, ok)
	assert.Equal(t, 9, start)
	assert.Equal(t, 12, end)

	start, end, ok = caps.PosName("snowman")
	assert.True(t, ok)
	assert.Equal(t, 9, start)
	assert.Equal(t, 12, end)
}

func TestIter(t *testing.T) {
	re := MustCompile(`\w+(\w)`)
	it := re.Iter("abc xyz")

	assert.True(t, it.Next())
	start, end := it.Match()
	assert.Equal(t, 0, start)
	assert.Equal(t, 3, end)

	assert.True(t, it.NextCaptures())
	caps := it.Captures()
	start, end, ok := caps.Pos(1)
	assert.True(t, ok)
	assert.Equal(t, 6, start)
	assert.Equal(t, 7, end)

	assert.False(t, it.Next())
	assert.False(t, it.NextCaptures())
}
