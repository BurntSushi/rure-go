package rure

import (
	"strings"
	"testing"
)

func BenchmarkLiteral(b *testing.B) {
	x := strings.Repeat("x", 50) + "y"
	b.StopTimer()
	re := MustCompile("y")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !re.IsMatch(x) {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkNotLiteral(b *testing.B) {
	x := strings.Repeat("x", 50) + "y"
	b.StopTimer()
	re := MustCompile(".y")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !re.IsMatch(x) {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkMatchClass(b *testing.B) {
	b.StopTimer()
	x := strings.Repeat("xxxx", 20) + "w"
	re := MustCompile("[abcdw]")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !re.IsMatch(x) {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkMatchClass_InRange(b *testing.B) {
	b.StopTimer()
	// 'b' is between 'a' and 'c', so the charclass
	// range checking is no help here.
	x := strings.Repeat("bbbb", 20) + "c"
	re := MustCompile("[ac]")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if !re.IsMatch(x) {
			b.Fatalf("no match!")
		}
	}
}

func BenchmarkAnchoredLiteralShortNonMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := MustCompile("^zbc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkAnchoredLiteralLongNonMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 15; i++ {
		x = append(x, x...)
	}
	re := MustCompile("^zbc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkAnchoredShortMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := MustCompile("^.bc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkAnchoredLongMatch(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 15; i++ {
		x = append(x, x...)
	}
	re := MustCompile("^.bc(d|e)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkOnePassShortA(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := MustCompile("^.bc(d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkNotOnePassShortA(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := MustCompile(".bc(d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkOnePassShortB(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := MustCompile("^.bc(?:d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkNotOnePassShortB(b *testing.B) {
	b.StopTimer()
	x := []byte("abcddddddeeeededd")
	re := MustCompile(".bc(?:d|e)*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkOnePassLongPrefix(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := MustCompile("^abcdefghijklmnopqrstuvwxyz.*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkOnePassLongNotPrefix(b *testing.B) {
	b.StopTimer()
	x := []byte("abcdefghijklmnopqrstuvwxyz")
	re := MustCompile("^.bcdefghijklmnopqrstuvwxyz.*$")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.IsMatchBytes(x)
	}
}

func BenchmarkMatchParallelShared(b *testing.B) {
	x := []byte("this is a long line that contains foo bar baz")
	re := MustCompile("foo (ba+r)? baz")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			re.IsMatchBytes(x)
		}
	})
}

var text []byte

func makeText(n int) []byte {
	if len(text) >= n {
		return text[:n]
	}
	text = make([]byte, n)
	x := ^uint32(0)
	for i := range text {
		x += x
		x ^= 1
		if int32(x) < 0 {
			x ^= 0x88888eef
		}
		if x%31 == 0 {
			text[i] = '\n'
		} else {
			text[i] = byte(x%(0x7E+1-0x20) + 0x20)
		}
	}
	return text
}

func benchmark(b *testing.B, re string, n int) {
	r := MustCompile(re)
	t := makeText(n)
	b.ResetTimer()
	b.SetBytes(int64(n))
	for i := 0; i < b.N; i++ {
		if r.IsMatchBytes(t) {
			b.Fatal("match!")
		}
	}
}

const (
	easy0  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ$"
	easy0i = "(?i)ABCDEFGHIJklmnopqrstuvwxyz$"
	easy1  = "A[AB]B[BC]C[CD]D[DE]E[EF]F[FG]G[GH]H[HI]I[IJ]J$"
	medium = "[XYZ]ABCDEFGHIJKLMNOPQRSTUVWXYZ$"
	hard   = "[ -~]*ABCDEFGHIJKLMNOPQRSTUVWXYZ$"
	hard1  = "ABCD|CDEF|EFGH|GHIJ|IJKL|KLMN|MNOP|OPQR|QRST|STUV|UVWX|WXYZ"
)

func BenchmarkMatchEasy0_32(b *testing.B)   { benchmark(b, easy0, 32<<0) }
func BenchmarkMatchEasy0_1K(b *testing.B)   { benchmark(b, easy0, 1<<10) }
func BenchmarkMatchEasy0_32K(b *testing.B)  { benchmark(b, easy0, 32<<10) }
func BenchmarkMatchEasy0_1M(b *testing.B)   { benchmark(b, easy0, 1<<20) }
func BenchmarkMatchEasy0_32M(b *testing.B)  { benchmark(b, easy0, 32<<20) }
func BenchmarkMatchEasy0i_32(b *testing.B)  { benchmark(b, easy0i, 32<<0) }
func BenchmarkMatchEasy0i_1K(b *testing.B)  { benchmark(b, easy0i, 1<<10) }
func BenchmarkMatchEasy0i_32K(b *testing.B) { benchmark(b, easy0i, 32<<10) }
func BenchmarkMatchEasy0i_1M(b *testing.B)  { benchmark(b, easy0i, 1<<20) }
func BenchmarkMatchEasy0i_32M(b *testing.B) { benchmark(b, easy0i, 32<<20) }
func BenchmarkMatchEasy1_32(b *testing.B)   { benchmark(b, easy1, 32<<0) }
func BenchmarkMatchEasy1_1K(b *testing.B)   { benchmark(b, easy1, 1<<10) }
func BenchmarkMatchEasy1_32K(b *testing.B)  { benchmark(b, easy1, 32<<10) }
func BenchmarkMatchEasy1_1M(b *testing.B)   { benchmark(b, easy1, 1<<20) }
func BenchmarkMatchEasy1_32M(b *testing.B)  { benchmark(b, easy1, 32<<20) }
func BenchmarkMatchMedium_32(b *testing.B)  { benchmark(b, medium, 32<<0) }
func BenchmarkMatchMedium_1K(b *testing.B)  { benchmark(b, medium, 1<<10) }
func BenchmarkMatchMedium_32K(b *testing.B) { benchmark(b, medium, 32<<10) }
func BenchmarkMatchMedium_1M(b *testing.B)  { benchmark(b, medium, 1<<20) }
func BenchmarkMatchMedium_32M(b *testing.B) { benchmark(b, medium, 32<<20) }
func BenchmarkMatchHard_32(b *testing.B)    { benchmark(b, hard, 32<<0) }
func BenchmarkMatchHard_1K(b *testing.B)    { benchmark(b, hard, 1<<10) }
func BenchmarkMatchHard_32K(b *testing.B)   { benchmark(b, hard, 32<<10) }
func BenchmarkMatchHard_1M(b *testing.B)    { benchmark(b, hard, 1<<20) }
func BenchmarkMatchHard_32M(b *testing.B)   { benchmark(b, hard, 32<<20) }
func BenchmarkMatchHard1_32(b *testing.B)   { benchmark(b, hard1, 32<<0) }
func BenchmarkMatchHard1_1K(b *testing.B)   { benchmark(b, hard1, 1<<10) }
func BenchmarkMatchHard1_32K(b *testing.B)  { benchmark(b, hard1, 32<<10) }
func BenchmarkMatchHard1_1M(b *testing.B)   { benchmark(b, hard1, 1<<20) }
func BenchmarkMatchHard1_32M(b *testing.B)  { benchmark(b, hard1, 32<<20) }
