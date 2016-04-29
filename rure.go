package rure

// #cgo LDFLAGS: -lrure
// #include <stdio.h>
// #include "rure.h"
//
// void rure_iter_collect(rure_iter *it,
//						  const uint8_t *haystack, size_t length,
//						  size_t **matches, size_t *num_matches)
// {
//     rure_match m = {0};
//     size_t len = 0;
//     size_t cap = 64;
//     *matches = (size_t *)malloc(cap * sizeof(size_t));
//     while (rure_iter_next(it, haystack, length, &m)) {
//	       if ((len * 2 + 1) >= cap) {
//	           cap *= 2;
//		       *matches = (size_t *)realloc(*matches, cap * sizeof(size_t));
//		   }
//		   (*matches)[len * 2 + 0] = m.start;
//		   (*matches)[len * 2 + 1] = m.end;
//         len++;
//     }
//     *num_matches = len;
// }
import "C"

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

const (
	FlagCaseI     = C.RURE_FLAG_CASEI
	FlagMulti     = C.RURE_FLAG_MULTI
	FlagDotNL     = C.RURE_FLAG_DOTNL
	FlagSwapGreed = C.RURE_FLAG_SWAP_GREED
	FlagSpace     = C.RURE_FLAG_SPACE
	FlagUnicode   = C.RURE_FLAG_UNICODE
	FlagDefault   = C.RURE_DEFAULT_FLAGS
)

type Regex struct {
	pattern string
	p       *C.rure
}

type Options struct {
	p *C.rure_options
}

type Captures struct {
	re *Regex
	p  *C.rure_captures
	ok bool
}

type Iter struct {
	re       *Regex
	p        *C.rure_iter
	haystack []byte
	match    C.rure_match
}

type Error struct {
	p *C.rure_error
}

func MustCompile(pattern string) *Regex {
	re, err := Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("regex.MustCompile failed on %s: %s", pattern, err))
	}
	return re
}

func Compile(pattern string) (*Regex, error) {
	return CompileOptions(pattern, FlagDefault, nil)
}

func CompileOptions(
	pattern string,
	flags uint32,
	options *Options,
) (*Regex, error) {
	re := &Regex{pattern: pattern}
	runtime.SetFinalizer(re, func(re *Regex) {
		if re.p != nil {
			C.rure_free(re.p)
			re.p = nil
		}
	})

	var optp *C.rure_options
	if options != nil {
		optp = options.p
	}
	err := newError()
	re.p = C.rure_compile(
		(*C.uint8_t)(&noCopyBytes(pattern)[0]),
		C.size_t(len(pattern)),
		FlagDefault,
		optp,
		err.p,
	)
	if re.p == nil {
		return nil, err
	}
	return re, nil
}

func (re *Regex) String() string {
	return re.pattern
}

func (re *Regex) IsMatch(text string) bool {
	return re.IsMatchBytes(noCopyBytes(text))
}

func (re *Regex) IsMatchBytes(text []byte) bool {
	return bool(C.rure_is_match(
		re.p, (*C.uint8_t)(&text[0]), C.size_t(len(text)), 0))
}

func (re *Regex) ShortestMatch(text string) (end int, ok bool) {
	return re.ShortestMatchBytes(noCopyBytes(text))
}

func (re *Regex) ShortestMatchBytes(text []byte) (end int, ok bool) {
	var cend C.size_t
	ok = bool(C.rure_shortest_match(
		re.p, (*C.uint8_t)(&text[0]), C.size_t(len(text)), 0, &cend))
	end = int(cend)
	return
}

func (re *Regex) Find(text string) (start, end int, ok bool) {
	return re.FindBytes(noCopyBytes(text))
}

func (re *Regex) FindBytes(text []byte) (start, end int, ok bool) {
	match := C.rure_match{}
	ok = bool(C.rure_find(
		re.p, (*C.uint8_t)(&text[0]), C.size_t(len(text)), 0, &match))
	if ok {
		start, end = int(match.start), int(match.end)
	}
	return
}

func (re *Regex) FindAll(text string) []int {
	return re.FindAllBytes(noCopyBytes(text))
}

func (re *Regex) FindAllBytes(text []byte) []int {
	it := C.rure_iter_new(re.p)
	defer C.rure_iter_free(it)

	haystack := (*C.uint8_t)(&text[0])
	len := C.size_t(len(text))
	nmatches := C.size_t(0)
	matches := (*C.size_t)(nil)
	defer C.free(unsafe.Pointer(matches))

	C.rure_iter_collect(it, haystack, len, &matches, &nmatches)
	if nmatches == 0 {
		return nil
	}
	matchesInts := make([]int, 2*nmatches)
	matchesArr := (*[1 << 30]C.size_t)(unsafe.Pointer(matches))
	for i := 0; i < int(nmatches); i++ {
		matchesInts[i*2] = int(matchesArr[i*2])
		matchesInts[i*2+1] = int(matchesArr[i*2+1])
	}
	return matchesInts
}

func (re *Regex) NewCaptures() *Captures {
	caps := &Captures{re: re, p: C.rure_captures_new(re.p)}
	runtime.SetFinalizer(caps, func(caps *Captures) {
		if caps.p != nil {
			C.rure_captures_free(caps.p)
			caps.p = nil
		}
	})
	return caps
}

func (re *Regex) Captures(caps *Captures, text string) bool {
	return re.CapturesBytes(caps, noCopyBytes(text))
}

func (re *Regex) CapturesBytes(caps *Captures, text []byte) bool {
	caps.ok = bool(C.rure_find_captures(
		re.p, (*C.uint8_t)(&text[0]), C.size_t(len(text)), 0, caps.p))
	return caps.ok
}

func (re *Regex) Iter(text string) *Iter {
	return re.IterBytes(noCopyBytes(text))
}

func (re *Regex) IterBytes(text []byte) *Iter {
	return newIter(re, text)
}

func NewOptions() *Options {
	opts := &Options{C.rure_options_new()}
	runtime.SetFinalizer(opts, func(opts *Options) {
		if opts.p != nil {
			C.rure_options_free(opts.p)
			opts.p = nil
		}
	})
	return opts
}

func (opts *Options) SetSizeLimit(limit int) {
	C.rure_options_size_limit(opts.p, C.size_t(limit))
}

func (opts *Options) SetDFASizeLimit(limit int) {
	C.rure_options_dfa_size_limit(opts.p, C.size_t(limit))
}

func (caps *Captures) IsMatch() bool {
	return caps.ok
}

func (caps *Captures) Group(i int) (start, end int, ok bool) {
	match := C.rure_match{}
	ok = bool(C.rure_captures_at(caps.p, C.size_t(i), &match))
	if ok {
		start, end = int(match.start), int(match.end)
	}
	return
}

func (caps *Captures) GroupName(name string) (start, end int, ok bool) {
	i := C.rure_capture_name_index(caps.re.p, C.CString(name))
	if i == -1 {
		return
	}
	return caps.Group(int(i))
}

func (caps *Captures) Len() int {
	return int(C.rure_captures_len(caps.p))
}

func newIter(re *Regex, haystack []byte) *Iter {
	it := &Iter{
		re:       re,
		haystack: haystack,
		p:        C.rure_iter_new(re.p),
	}
	runtime.SetFinalizer(it, func(it *Iter) {
		if it.p != nil {
			C.rure_iter_free(it.p)
		}
	})
	return it
}

func (it *Iter) Next(caps *Captures) bool {
	haystack := (*C.uint8_t)(&it.haystack[0])
	len := C.size_t(len(it.haystack))

	if caps == nil {
		return bool(C.rure_iter_next(it.p, haystack, len, &it.match))
	}
	caps.ok = bool(C.rure_iter_next_captures(it.p, haystack, len, caps.p))
	C.rure_captures_at(caps.p, 0, &it.match)
	return caps.ok
}

func (it *Iter) Match() (start, end int) {
	return int(it.match.start), int(it.match.end)
}

func newError() *Error {
	err := &Error{C.rure_error_new()}
	runtime.SetFinalizer(err, func(err *Error) {
		if err.p != nil {
			C.rure_error_free(err.p)
			err.p = nil
		}
	})
	return err
}

func (err *Error) Error() string {
	return C.GoString(C.rure_error_message(err.p))
}

// Converts a string to a []byte without allocating.
//
// This is very dangerous and must be handled with care. In particular, the
// input given must be kept alive for the duration of the return value.
func noCopyBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Copies a Go byte slice into C allocated memory.
//
// The caller is responsible for freeing the memory returned.
func cbytes(b []byte) unsafe.Pointer {
	// take from go:src/cmd/cgo/out.go
	p := C.malloc(C.size_t(len(b)))
	pp := (*[1 << 30]byte)(p)
	copy(pp[:], b)
	return p
}
