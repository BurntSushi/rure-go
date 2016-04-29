package rure

import "fmt"

// ExampleRegex shows how to compile a Regex, iterate over successive matches
// and extract a capturing group using its name.
func ExampleRegex() {
	re := MustCompile(`\w*(?P<last>\w)`)
	haystack := "foo bar baz quux"
	it := re.Iter(haystack)
	caps := re.NewCaptures()
	for it.Next(caps) {
		// Print the last letter of each word matched.
		start, end, _ := caps.GroupName("last")
		fmt.Println(haystack[start:end])
	}
	// Output:
	// o
	// r
	// z
	// x
}

// ExampleFindBytes shows how to match bytes that aren't UTF-8.
func ExampleRegex_FindBytes() {
	re := MustCompile(`(?-u)\xFF`)
	fmt.Println(re.FindBytes([]byte{0xFF}))
	// Output:
	// 0 1 true
}
