Go bindings to RUst's REgex engine
==================================
This package provides cgo bindings to
[Rust's regex engine](https://github.com/rust-lang/regex)
using its
[C API](https://github.com/rust-lang/regex/tree/master/regex-capi).

Dual-licensed under MIT or the [UNLICENSE](http://unlicense.org).


### Documentation

[godoc.org/github.com/BurntSushi/rure-go](http://godoc.org/github.com/BurntSushi/rure-go)

The primary documentation for the Rust library, including a definition of the
syntax, can be found here:
https://docs.rs/regex/#syntax


### Install

You'll need to [install Rust](https://www.rust-lang.org/downloads.html) (you'll
need at least Rust 1.9, which is the current beta release) and have a Go
compiler handy. To run tests for `rure-go`, we'll need to compile Rust's regex
library and then tell the Go compiler where to find it. These commands should
do it:

```
$ git clone git://github.com/rust-lang/regex
$ cargo build --release --manifest-path ./regex/regex-capi/Cargo.toml
$ export CGO_LDFLAGS="-L$(pwd)/regex/target/release"
$ export LD_LIBRARY_PATH="$(pwd)/regex/target/release"
$ go get -t github.com/BurntSushi/rure-go
$ go test github.com/BurntSushi/rure-go
```

And to run benchmarks:

```
$ go test github.com/BurntSushi/rure-go/bench/rust -cpu 1 -run / -bench .
```

Replace `rust` with `go` in the package path to run the same set of benchmarks
using Go's `regexp` package.


### Example usage

This shows how to compile a regex, iterate over successive matches and extract
a capturing group using its name:

```
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
```


### Motivation

I wrote these bindings primarily as a test case for the
[C API of Rust's regex library](https://github.com/rust-lang/regex/tree/master/regex-capi).
In particular, I wanted to be sure that it was feasible to write low overhead
bindings to another language. For the most part, I think that was a success.

Secondarily, Rust's regex engine is pretty fast and also provides the same
algorithmic guarantees that Go's regex engine provides. Therefore, it may prove
useful as an alternative under some limited circumstances. First, you must be
willing to tolerate cgo. Second, because of cgo function call overhead, your
haystacks probably need to be a bit bigger to realize big performance wins. For
example, if your use case is matching a regex a lot on very small strings, then
cgo will probably prevent you from doing it faster with `rure`. (You could,
however, write your matching loop as a C helper function, but now we're just
digging ourselves in deeper.)


### Benchmarks

These are the benchmarks as they are defined in Go's `regexp` package. Two
important comments:

1. Rust's regex engine hits a performance ceiling because of cgo function call
   overhead. As a result, Go is faster on some benchmarks.
2. Several of the benchmarks are dramatically distorted because of
   optimizations in Rust's regex engine. In fact, none of these benchmarks give
   an accurate depiction of Rust's regex engine overall speed. For example,
   the Hard1 benchmark doesn't even use the regex engine at all, because Rust's
   library will detect it as a pure alternation of literals and fall back to
   Aho-Corasick. The other {easy,medium,hard} benchmarks run the regex engine
   in reverse from the end of the string since the regexes are anchored at the
   end. As a result, the time it takes to execute the search doesn't vary with
   the length of the haystack.

```
$ benchcmp ~/clones/go/src/regexp/regex-go ./regex-rust
benchmark                                 old ns/op       new ns/op     delta
BenchmarkLiteral                          140             229           +63.57%
BenchmarkNotLiteral                       2099            328           -84.37%
BenchmarkMatchClass                       3252            295           -90.93%
BenchmarkMatchClass_InRange               2957            243           -91.78%
BenchmarkAnchoredLiteralShortNonMatch     118             222           +88.14%
BenchmarkAnchoredLiteralLongNonMatch      133             223           +67.67%
BenchmarkAnchoredShortMatch               189             226           +19.58%
BenchmarkAnchoredLongMatch                364             226           -37.91%
BenchmarkOnePassShortA                    808             255           -68.44%
BenchmarkNotOnePassShortA                 867             254           -70.70%
BenchmarkOnePassShortB                    571             256           -55.17%
BenchmarkNotOnePassShortB                 621             255           -58.94%
BenchmarkOnePassLongPrefix                131             272           +107.63%
BenchmarkOnePassLongNotPrefix             485             272           -43.92%
BenchmarkMatchParallelShared              332             279           -15.96%
BenchmarkMatchEasy0_32                    104             222           +113.46%
BenchmarkMatchEasy0_1K                    648             222           -65.74%
BenchmarkMatchEasy0_32K                   11524           222           -98.07%
BenchmarkMatchEasy0_1M                    423966          222           -99.95%
BenchmarkMatchEasy0_32M                   14098625        213           -100.00%
BenchmarkMatchEasy0i_32                   2323            228           -90.19%
BenchmarkMatchEasy0i_1K                   68757           228           -99.67%
BenchmarkMatchEasy0i_32K                  2711922         228           -99.99%
BenchmarkMatchEasy0i_1M                   86969502        228           -100.00%
BenchmarkMatchEasy0i_32M                  2775970393      233           -100.00%
BenchmarkMatchEasy1_32                    86.0            229           +166.28%
BenchmarkMatchEasy1_1K                    1455            229           -84.26%
BenchmarkMatchEasy1_32K                   52058           229           -99.56%
BenchmarkMatchEasy1_1M                    1825699         229           -99.99%
BenchmarkMatchEasy1_32M                   58742205        211           -100.00%
BenchmarkMatchMedium_32                   1361            242           -82.22%
BenchmarkMatchMedium_1K                   40203           242           -99.40%
BenchmarkMatchMedium_32K                  1713513         242           -99.99%
BenchmarkMatchMedium_1M                   54703006        242           -100.00%
BenchmarkMatchMedium_32M                  1755415300      214           -100.00%
BenchmarkMatchHard_32                     2183            229           -89.51%
BenchmarkMatchHard_1K                     66067           229           -99.65%
BenchmarkMatchHard_32K                    3113396         229           -99.99%
BenchmarkMatchHard_1M                     99616880        229           -100.00%
BenchmarkMatchHard_32M                    3189585261      214           -100.00%
BenchmarkMatchHard1_32                    12668           272           -97.85%
BenchmarkMatchHard1_1K                    395644          2333          -99.41%
BenchmarkMatchHard1_32K                   16404742        68659         -99.58%
BenchmarkMatchHard1_1M                    525035229       2186555       -99.58%
BenchmarkMatchHard1_32M                   16667981723     70463061      -99.58%

benchmark                    old MB/s     new MB/s         speedup
BenchmarkMatchEasy0_32       306.08       143.69           0.47x
BenchmarkMatchEasy0_1K       1578.82      4605.20          2.92x
BenchmarkMatchEasy0_32K      2843.36      147213.66        51.77x
BenchmarkMatchEasy0_1M       2473.25      4719689.39       1908.29x
BenchmarkMatchEasy0_32M      2379.98      156878555.03     65915.91x
BenchmarkMatchEasy0i_32      13.77        139.75           10.15x
BenchmarkMatchEasy0i_1K      14.89        4475.04          300.54x
BenchmarkMatchEasy0i_32K     12.08        143235.97        11857.28x
BenchmarkMatchEasy0i_1M      12.06        4582492.43       379974.50x
BenchmarkMatchEasy0i_32M     12.09        143626597.90     11879784.77x
BenchmarkMatchEasy1_32       372.19       139.49           0.37x
BenchmarkMatchEasy1_1K       703.69       4467.93          6.35x
BenchmarkMatchEasy1_32K      629.44       143010.92        227.20x
BenchmarkMatchEasy1_1M       574.34       4575198.42       7966.01x
BenchmarkMatchEasy1_32M      571.22       158747658.13     277909.84x
BenchmarkMatchMedium_32      23.50        132.01           5.62x
BenchmarkMatchMedium_1K      25.47        4231.03          166.12x
BenchmarkMatchMedium_32K     19.12        135348.39        7078.89x
BenchmarkMatchMedium_1M      19.17        4331561.73       225955.23x
BenchmarkMatchMedium_32M     19.11        156575833.69     8193397.89x
BenchmarkMatchHard_32        14.66        139.50           9.52x
BenchmarkMatchHard_1K        15.50        4466.01          288.13x
BenchmarkMatchHard_32K       10.52        142758.39        13570.19x
BenchmarkMatchHard_1M        10.53        4574748.55       434449.06x
BenchmarkMatchHard_32M       10.52        156417047.69     14868540.65x
BenchmarkMatchHard1_32       2.53         117.23           46.34x
BenchmarkMatchHard1_1K       2.59         438.78           169.41x
BenchmarkMatchHard1_32K      2.00         477.25           238.62x
BenchmarkMatchHard1_1M       2.00         479.56           239.78x
BenchmarkMatchHard1_32M      2.01         476.20           236.92x
```

Here are the benchmarks ported from Rust's regex library. These benchmarks do
a better job of characterizing the performance difference. In particular, they
use `FindAll()` on a single `0.5MB` haystack, which iterates over successive
matches in C to avoid cgo function call overhead.

```
benchmark                               old MB/s     new MB/s     speedup
BenchmarkSherlock                       3805.25      7861.29      2.07x
BenchmarkHolmes                         1339.82      9515.28      7.10x
BenchmarkSherlockHolmes                 2900.11      14011.42     4.83x
BenchmarkSherlockNoCase                 9.76         509.87       52.24x
BenchmarkHolmesNoCase                   10.66        534.47       50.14x
BenchmarkSherlockHolmesNoCase           9.84         516.71       52.51x
BenchmarkNameWhitespace                 2591.41      6870.12      2.65x
BenchmarkNameAlt1                       1556.22      12787.94     8.22x
BenchmarkNameAlt2                       10.20        2789.56      273.49x
BenchmarkNameAlt3                       2.93         485.43       165.68x
BenchmarkNameAlt3NoCase                 1.56         429.32       275.21x
BenchmarkNameAlt4                       10.17        2311.07      227.24x
BenchmarkNameAlt4NoCase                 5.31         434.40       81.81x
BenchmarkNameAlt5                       7.08         1735.48      245.12x
BenchmarkNameAlt5NoCase                 3.55         443.53       124.94x
BenchmarkNoMatchUncommon                15335.81     21447.14     1.40x
BenchmarkNoMatchCommon                  487.44       21294.92     43.69x
BenchmarkNoMatchReallyCommon            468.76       1622.89      3.46x
BenchmarkTheLower                       100.82       755.81       7.50x
BenchmarkTheUpper                       1219.33      8116.70      6.66x
BenchmarkTheNoCase                      9.68         324.15       33.49x
BenchmarkTheWhitespace                  59.02        442.80       7.50x
BenchmarkEverythingGreedy               8.86         179.60       20.27x
BenchmarkEverythingGreedyNL             10.75        472.34       43.94x
BenchmarkLetters                        1.72         15.70        9.13x
BenchmarkLettersUpper                   13.99        247.74       17.71x
BenchmarkLettersLower                   1.64         16.43        10.02x
BenchmarkWords                          4.19         48.45        11.56x
BenchmarkBeforeHolmes                   8.58         517.72       60.34x
BenchmarkHolmesCocharWatson             10.73        2634.69      245.54x
BenchmarkHolmesCowordWatson             2.15         883.32       410.85x
BenchmarkQuotes                         12.49        1002.75      80.28x
BenchmarkLineBoundarySherlockHolmes     12.98        534.98       41.22x
BenchmarkWordEndingN                    8.82         288.50       32.71x
BenchmarkRepeatedClassNegation          4.68         1.09         0.23x
BenchmarkIngSuffix                      10.77        413.77       38.42x
BenchmarkIngSuffixLimitedSpace          6.96         437.96       62.93x
```

The principal reason for the performance difference is that Rust's regex engine
(like RE2) has a lazy DFA.
[There is ongoing work](https://github.com/golang/go/issues/11646)
to add a lazy DFA to Go's `regexp` package.

For fun, and since Rust's regex engine, RE2 (in C++) and Go's regexp library
are very much related and share many of the same implementation details,
here's a comparison between Rust's regex engine and RE2 using Rust's benchmark
harness:

```
name                                     re2 ns/iter           rust ns/iter            diff ns/iter   diff %
sherlock::before_holmes                  1,443,662 (412 MB/s)  1,142,847 (520 MB/s)        -300,815  -20.84%
sherlock::everything_greedy              9,310,589 (63 MB/s)   2,516,121 (236 MB/s)      -6,794,468  -72.98%
sherlock::everything_greedy_nl           2,491,676 (238 MB/s)  1,204,517 (493 MB/s)      -1,287,159  -51.66%
sherlock::holmes_cochar_watson           1,175,610 (506 MB/s)  222,959 (2,668 MB/s)        -952,651  -81.03%
sherlock::holmes_coword_watson           1,347,865 (441 MB/s)  638,450 (931 MB/s)          -709,415  -52.63%
sherlock::ing_suffix                     3,028,970 (196 MB/s)  1,361,756 (436 MB/s)      -1,667,214  -55.04%
sherlock::ing_suffix_limited_space       1,957,296 (303 MB/s)  1,302,217 (456 MB/s)        -655,079  -33.47%
sherlock::letters                        111,006,069 (5 MB/s)  25,340,930 (23 MB/s)     -85,665,139  -77.17%
sherlock::letters_lower                  107,183,616 (5 MB/s)  24,525,566 (24 MB/s)     -82,658,050  -77.12%
sherlock::letters_upper                  4,709,166 (126 MB/s)  2,082,277 (285 MB/s)      -2,626,889  -55.78%
sherlock::line_boundary_sherlock_holmes  2,552,194 (233 MB/s)  1,114,865 (533 MB/s)      -1,437,329  -56.32%
sherlock::name_alt1                      77,602 (7,666 MB/s)   37,361 (15,923 MB/s)         -40,241  -51.86%
sherlock::name_alt2                      1,320,543 (450 MB/s)  191,247 (3,110 MB/s)      -1,129,296  -85.52%
sherlock::name_alt3                      1,439,985 (413 MB/s)  1,280,912 (464 MB/s)        -159,073  -11.05%
sherlock::name_alt3_nocase               2,756,402 (215 MB/s)  1,339,666 (444 MB/s)      -1,416,736  -51.40%
sherlock::name_alt4                      1,362,748 (436 MB/s)  231,035 (2,575 MB/s)      -1,131,713  -83.05%
sherlock::name_alt4_nocase               2,025,273 (293 MB/s)  1,313,315 (453 MB/s)        -711,958  -35.15%
sherlock::name_alt5                      1,347,991 (441 MB/s)  322,420 (1,845 MB/s)      -1,025,571  -76.08%
sherlock::name_alt5_nocase               2,115,018 (281 MB/s)  1,314,237 (452 MB/s)        -800,781  -37.86%
sherlock::name_holmes                    166,170 (3,580 MB/s)  47,617 (12,494 MB/s)        -118,553  -71.34%
sherlock::name_holmes_nocase             1,647,206 (361 MB/s)  1,114,830 (533 MB/s)        -532,376  -32.32%
sherlock::name_sherlock                  58,758 (10,125 MB/s)  69,917 (8,509 MB/s)           11,159   18.99%
sherlock::name_sherlock_holmes           60,607 (9,816 MB/s)   36,615 (16,248 MB/s)         -23,992  -39.59%
sherlock::name_sherlock_holmes_nocase    1,540,731 (386 MB/s)  1,160,475 (512 MB/s)        -380,256  -24.68%
sherlock::name_sherlock_nocase           1,536,512 (387 MB/s)  1,161,208 (512 MB/s)        -375,304  -24.43%
sherlock::name_whitespace                62,172 (9,569 MB/s)   78,795 (7,550 MB/s)           16,623   26.74%
sherlock::no_match_common                440,514 (1,350 MB/s)  25,933 (22,941 MB/s)        -414,581  -94.11%
sherlock::no_match_really_common         440,342 (1,351 MB/s)  363,675 (1,635 MB/s)         -76,667  -17.41%
sherlock::no_match_uncommon              23,726 (25,075 MB/s)  25,809 (23,051 MB/s)           2,083    8.78%
sherlock::quotes                         1,393,911 (426 MB/s)  557,791 (1,066 MB/s)        -836,120  -59.98%
sherlock::the_lower                      2,504,957 (237 MB/s)  620,410 (958 MB/s)        -1,884,547  -75.23%
sherlock::the_nocase                     3,695,971 (160 MB/s)  1,681,153 (353 MB/s)      -2,014,818  -54.51%
sherlock::the_upper                      232,581 (2,557 MB/s)  50,412 (11,801 MB/s)        -182,169  -78.32%
sherlock::the_whitespace                 2,279,605 (260 MB/s)  1,121,804 (530 MB/s)      -1,157,801  -50.79%
sherlock::word_ending_n                  3,330,263 (178 MB/s)  2,015,397 (295 MB/s)      -1,314,866  -39.48%
sherlock::words                          30,546,351 (19 MB/s)  9,588,968 (62 MB/s)      -20,957,383  -68.61%
```

Some possible explanations for the performance difference:

1. Rust's regex engine does quite a bit better with prefix literal detection.
   Instead of limiting itself to just a single byte (like RE2), it will expand
   alternates and character classes (within some limit) to come up with
   literals. It will then use `memchr`, Boyer-Moore or Aho-Corasick before
   diving into the regex engine.
2. Rust's inner DFA loop is explicitly unrolled a few times, which decreases
   the number of instructions it has to execute. (RE2's DFA loop is not
   unrolled.)

There's probably more, but my experience suggests the above are the most
significant.


### Bugs

Rust's regex engine uses the equivalent of C's `size_t` type to represent match
offsets into a haystack. This FFI wrapper converts such offsets to Go's `int`
type, which is not guaranteed to have the same size as a `size_t`, and of
course, `int` is signed while `size_t` is not.

It's not clear what the right answer is here. Today, the conversions are
unchecked, which means that callers will get wrong answers if the size of the
haystack exceeds the size of Go's `int` type.
