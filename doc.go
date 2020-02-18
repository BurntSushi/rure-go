/*
Package rure contains Go bindings to Rust's regex library.

Like Go's regexp package in the standard library, rure will complete all
searches in linear time with respect to the search text.

The syntax has very few differences from the syntax supported by
Go's regexp package. Nevertheless, the syntax is documented here:
https://docs.rs/regex/#syntax The differences include, but are not limited to:
word boundary (\b) assertions are Unicode aware by default, Unicode support can
be toggled inside the expression with the u flag, and an x flag is supported to
permit writing patterns with insignificant whitespace.

Text encoding of patterns

All regular expression patterns must be valid UTF-8. Compilation will fail if
a pattern is invalid UTF-8. In order to match a sequence of bytes that is not
valid UTF-8, you can disable Unicode support and use hexadecimal escape
sequences. For example, to match the single byte \xFF (which is invalid UTF-8),
one can do:

	re := MustCompile(`(?-u)\xFF`)
	fmt.Println(re.FindBytes([]byte{0xFF}))
	// Output:
	// 0 1 true

If the pattern omitted (?-u), then \xFF would correspond to the Unicode
codepoint U+00FF, which has a UTF-8 encoding of \xCE\xBF, which of course
does not match the single byte \xFF.

Text encoding of haystacks

To a first approximation, haystacks should be UTF-8. In fact, UTF-8 (and, one
supposes, ASCII) is the only well defined text encoding supported by this
library. It is impossible to match UTF-16, UTF-32 or any other encoding without
first transcoding it to UTF-8.

With that said, haystacks do not need to be valid UTF-8, and if they aren't
valid UTF-8, no performance penalty is paid. Whether invalid UTF-8 is matched
or not depends on the regular expression. For example, with the FlagUnicode
flag enabled, the regex . is guaranteed to match a single UTF-8 encoding of a
Unicode codepoint (sans LF). In particular, it will not match invalid UTF-8
such as \xFF, nor will it match surrogate codepoints or "alternate" (i.e.,
non-minimal) encodings of codepoints. However, with the FlagUnicode flag
disabled, the regex . will match any single arbitrary byte (sans LF), including
\xFF.

This provides a useful invariant: wherever FlagUnicode is set, the
corresponding regex is guaranteed to match valid UTF-8. Invalid UTF-8 will
always prevent a match from happening when the flag is set. Since flags can be
toggled in the regular expression itself, this allows one to pick and choose
which parts of the regular expression must match UTF-8 or not.

Some good advice is to always enable the FlagUnicode flag (which is enabled
when using MustCompile) and selectively disable the flag when one wants to
match arbitrary bytes. The flag can be disabled in a regular expression with
(?-u).

Performance tips

When matching text, prefer methods in this order: IsMatch, Find, Captures.
e.g., Using Captures without inspecting a non-zero capture group will cause the
regex engine to do wasted work.
*/
package rure
