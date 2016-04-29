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

func TestShortestMatch(t *testing.T) {
	re := MustCompile(`a+`)
	end, ok := re.ShortestMatch("aaaaa")
	assert.True(t, ok)
	assert.Equal(t, 1, end)
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
	caps := re.NewCaptures()
	assert.Equal(t, 3, caps.Len())

	ok := re.Captures(caps, "snowman: ☃")
	assert.True(t, ok)
	assert.NotNil(t, caps)

	start, end, ok := caps.Group(2)
	assert.True(t, ok)
	assert.Equal(t, 9, start)
	assert.Equal(t, 12, end)

	start, end, ok = caps.GroupName("snowman")
	assert.True(t, ok)
	assert.Equal(t, 9, start)
	assert.Equal(t, 12, end)
}

func TestIter(t *testing.T) {
	re := MustCompile(`\w+(\w)`)
	it := re.Iter("abc xyz")

	ok := it.Next(nil)
	assert.True(t, ok)
	start, end := it.Match()
	assert.Equal(t, 0, start)
	assert.Equal(t, 3, end)

	caps := re.NewCaptures()
	ok = it.Next(caps)
	assert.True(t, ok)
	start, end = it.Match()
	assert.Equal(t, 4, start)
	assert.Equal(t, 7, end)

	start, end, ok = caps.Group(1)
	assert.True(t, ok)
	assert.Equal(t, 6, start)
	assert.Equal(t, 7, end)

	assert.False(t, it.Next(nil))
}

func TestFindAll(t *testing.T) {
	re := MustCompile(`\w+(\w)`)
	matches := re.FindAll("abc xyz")
	assert.Equal(t, []int{0, 3, 4, 7}, matches)
}
