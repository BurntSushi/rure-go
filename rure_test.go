package rure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	re, err := Compile(`(`)
	require.Nil(t, re)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unclosed group")
}

func TestIsMatch(t *testing.T) {
	re := MustCompile(`\p{So}`)
	require.True(t, re.IsMatch("snowman: ☃"))
}

func TestShortestMatch(t *testing.T) {
	re := MustCompile(`a+`)
	end, ok := re.ShortestMatch("aaaaa")
	require.True(t, ok)
	require.Equal(t, 1, end)
}

func TestFind(t *testing.T) {
	re := MustCompile(`\p{So}`)
	start, end, ok := re.Find("snowman: ☃")
	require.True(t, ok)
	require.Equal(t, 9, start)
	require.Equal(t, 12, end)
}

func TestCaptures(t *testing.T) {
	re := MustCompile(`.(.*(?P<snowman>\p{So}))$`)
	caps := re.NewCaptures()
	require.Equal(t, 3, caps.Len())

	ok := re.Captures(caps, "snowman: ☃")
	require.True(t, ok)
	require.NotNil(t, caps)

	start, end, ok := caps.Group(2)
	require.True(t, ok)
	require.Equal(t, 9, start)
	require.Equal(t, 12, end)

	start, end, ok = caps.GroupName("snowman")
	require.True(t, ok)
	require.Equal(t, 9, start)
	require.Equal(t, 12, end)
}

func TestIter(t *testing.T) {
	re := MustCompile(`\w+(\w)`)
	it := re.Iter("abc xyz")

	ok := it.Next(nil)
	require.True(t, ok)
	start, end := it.Match()
	require.Equal(t, 0, start)
	require.Equal(t, 3, end)

	caps := re.NewCaptures()
	ok = it.Next(caps)
	require.True(t, ok)
	start, end = it.Match()
	require.Equal(t, 4, start)
	require.Equal(t, 7, end)

	start, end, ok = caps.Group(1)
	require.True(t, ok)
	require.Equal(t, 6, start)
	require.Equal(t, 7, end)

	require.False(t, it.Next(nil))
}

func TestFindAll(t *testing.T) {
	re := MustCompile(`\w+(\w)`)
	matches := re.FindAll("abc xyz")
	require.Equal(t, []int{0, 3, 4, 7}, matches)
}

func TestAt(t *testing.T) {
	re := MustCompile(`\bbar`)
	haystack := "foobar"
	require.True(t, re.IsMatch(haystack[3:]))
	require.False(t, re.IsMatchAt(haystack, 3))
}

func TestCaptureNames(t *testing.T) {
	re := MustCompile(`(?P<foo>zzz)(zzz)(?:zzz)(?P<bar>zzz)`)
	require.Equal(t, []string{"", "foo", "", "bar"}, re.CaptureNames())
}
