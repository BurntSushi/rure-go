package rure_bench_rust

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/BurntSushi/rure-go"
)

func BenchmarkSherlock(b *testing.B) {
	bench(b, `Sherlock`, 97)
}

func BenchmarkHolmes(b *testing.B) {
	bench(b, `Holmes`, 461)
}

func BenchmarkSherlockHolmes(b *testing.B) {
	bench(b, `Sherlock Holmes`, 91)
}

func BenchmarkSherlockNoCase(b *testing.B) {
	bench(b, `(?i)Sherlock`, 102)
}

func BenchmarkHolmesNoCase(b *testing.B) {
	bench(b, `(?i)Holmes`, 467)
}

func BenchmarkSherlockHolmesNoCase(b *testing.B) {
	bench(b, `(?i)Sherlock Holmes`, 96)
}

func BenchmarkNameWhitespace(b *testing.B) {
	bench(b, `Sherlock\s+Holmes`, 97)
}

func BenchmarkNameAlt1(b *testing.B) {
	bench(b, `Sherlock|Street`, 158)
}

func BenchmarkNameAlt2(b *testing.B) {
	bench(b, `Sherlock|Holmes`, 558)
}

func BenchmarkNameAlt3(b *testing.B) {
	bench(b, `Sherlock|Holmes|Watson|Irene|Adler|John|Baker`, 740)
}

func BenchmarkNameAlt3NoCase(b *testing.B) {
	bench(b, `(?i)Sherlock|Holmes|Watson|Irene|Adler|John|Baker`, 753)
}

func BenchmarkNameAlt4(b *testing.B) {
	bench(b, `Sher[a-z]+|Hol[a-z]+`, 582)
}

func BenchmarkNameAlt4NoCase(b *testing.B) {
	bench(b, `(?i)Sher[a-z]+|Hol[a-z]+`, 697)
}

func BenchmarkNameAlt5(b *testing.B) {
	bench(b, `Sherlock|Holmes|Watson`, 639)
}

func BenchmarkNameAlt5NoCase(b *testing.B) {
	bench(b, `(?i)Sherlock|Holmes|Watson`, 650)
}

func BenchmarkNoMatchUncommon(b *testing.B) {
	bench(b, `zqj`, 0)
}

func BenchmarkNoMatchCommon(b *testing.B) {
	bench(b, `aqj`, 0)
}

func BenchmarkNoMatchReallyCommon(b *testing.B) {
	bench(b, `aei`, 0)
}

func BenchmarkTheLower(b *testing.B) {
	bench(b, `the`, 7218)
}

func BenchmarkTheUpper(b *testing.B) {
	bench(b, `The`, 741)
}

func BenchmarkTheNoCase(b *testing.B) {
	bench(b, `(?i)the`, 7987)
}

func BenchmarkTheWhitespace(b *testing.B) {
	bench(b, `the\s+\w+`, 5410)
}

func BenchmarkEverythingGreedy(b *testing.B) {
	bench(b, `.*`, 13053)
}

func BenchmarkEverythingGreedyNL(b *testing.B) {
	bench(b, `(?s).*`, 1)
}

func BenchmarkLetters(b *testing.B) {
	bench(b, `\p{L}`, 447160)
}

func BenchmarkLettersUpper(b *testing.B) {
	bench(b, `\p{Lu}`, 14180)
}

func BenchmarkLettersLower(b *testing.B) {
	bench(b, `\p{Ll}`, 432980)
}

func BenchmarkWords(b *testing.B) {
	bench(b, `\w+`, 109214)
}

func BenchmarkBeforeHolmes(b *testing.B) {
	bench(b, `\w+\s+Holmes`, 319)
}

func BenchmarkHolmesCocharWatson(b *testing.B) {
	bench(b, `Holmes.{0,25}Watson|Watson.{0,25}Holmes`, 7)
}

func BenchmarkHolmesCowordWatson(b *testing.B) {
	bench(
		b,
		`Holmes(?:\s*.+\s*){0,10}Watson|Watson(?:\s*.+\s*){0,10}Holmes`,
		51)
}

func BenchmarkQuotes(b *testing.B) {
	bench(b, `["'][^"']{0,30}[?!.]["']`, 767)
}

func BenchmarkLineBoundarySherlockHolmes(b *testing.B) {
	bench(b, `(?m)^Sherlock Holmes|Sherlock Holmes$`, 34)
}

func BenchmarkWordEndingN(b *testing.B) {
	// Rust's regex engine uses Unicode word boundaries. Force to ASCII.
	bench(b, `(?-u:\b)\w+n(?-u:\b)`, 8366)
}

func BenchmarkRepeatedClassNegation(b *testing.B) {
	bench(b, `[a-q][^u-z]{13}x`, 142)
}

func BenchmarkIngSuffix(b *testing.B) {
	bench(b, `[a-zA-Z]+ing`, 2824)
}

func BenchmarkIngSuffixLimitedSpace(b *testing.B) {
	bench(b, `\s[a-zA-Z]{0,12}ing\s`, 2081)
}

func bench(b *testing.B, pattern string, expectCount int) {
	haystack := readFile("../../testdata/sherlock.txt")
	benchCount(b, haystack, pattern, expectCount)
}

func benchCount(
	b *testing.B,
	haystack []byte,
	pattern string,
	expectCount int,
) {
	re := rure.MustCompile(pattern)
	b.ResetTimer()
	b.SetBytes(int64(len(haystack)))
	for i := 0; i < b.N; i++ {
		count := len(re.FindAllBytes(haystack)) / 2
		if expectCount != count {
			panic(fmt.Sprintf(
				"expected %d matches but got %d", expectCount, count))
		}
	}
}

func readFile(fp string) []byte {
	f, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		panic(err)
	}
	return data
}
