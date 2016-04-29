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
