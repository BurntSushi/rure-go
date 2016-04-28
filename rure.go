package rure

// #cgo LDFLAGS: -lrure
// #include "rure.h"
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
}

type Iter struct {
	re    *Regex
	p     *C.rure_iter
	hayp  *C.uint8_t
	ok    bool
	match C.rure_match
	caps  *Captures
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

	err := newError()
	re.p = C.rure_compile(
		(*C.uint8_t)(&noCopyBytes(pattern)[0]),
		C.size_t(len(pattern)),
		FlagDefault,
		nil,
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

func (re *Regex) Captures(text string) *Captures {
	return re.CapturesBytes(noCopyBytes(text))
}

func (re *Regex) CapturesBytes(text []byte) *Captures {
	caps := newCaptures(re)
	ok := bool(C.rure_find_captures(
		re.p, (*C.uint8_t)(&text[0]), C.size_t(len(text)), 0, caps.p))
	if !ok {
		return nil
	}
	return caps
}

func (re *Regex) Iter(text string) *Iter {
	return re.IterBytes(noCopyBytes(text))
}

func (re *Regex) IterBytes(text []byte) *Iter {
	return newIter(re, text)
}

func newOptions() *Options {
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

func newCaptures(re *Regex) *Captures {
	caps := &Captures{re: re, p: C.rure_captures_new(re.p)}
	runtime.SetFinalizer(caps, func(caps *Captures) {
		if caps.p != nil {
			C.rure_captures_free(caps.p)
			caps.p = nil
		}
	})
	return caps
}

func (caps *Captures) Pos(i int) (start, end int, ok bool) {
	match := C.rure_match{}
	ok = bool(C.rure_captures_at(caps.p, C.size_t(i), &match))
	if ok {
		start, end = int(match.start), int(match.end)
	}
	return
}

func (caps *Captures) PosName(name string) (start, end int, ok bool) {
	i := C.rure_capture_name_index(caps.re.p, C.CString(name))
	if i == -1 {
		return
	}
	return caps.Pos(int(i))
}

func newIter(re *Regex, haystack []byte) *Iter {
	hayp := (*C.uint8_t)(cbytes(haystack))
	it := &Iter{
		re:   re,
		p:    C.rure_iter_new(re.p, hayp, C.size_t(len(haystack))),
		hayp: hayp,
	}
	runtime.SetFinalizer(it, func(it *Iter) {
		if it.p != nil {
			C.rure_iter_free(it.p)
			C.free(it.hayp)
		}
	})
	return it
}

func (it *Iter) Next() bool {
	it.ok = false
	it.caps = nil
	it.ok = bool(C.rure_iter_next(it.p, &it.match))
	return it.ok
}

func (it *Iter) NextCaptures() bool {
	it.ok = false
	it.caps = newCaptures(it.re)
	it.ok = bool(C.rure_iter_next_captures(it.p, it.caps.p))
	return it.ok
}

func (it *Iter) Match() (start, end int) {
	if !it.ok || it.caps != nil {
		panic("Next must return true before calling Match")
	}
	return int(it.match.start), int(it.match.end)
}

func (it *Iter) Captures() *Captures {
	if !it.ok || it.caps == nil {
		panic("NextCaptures must return true before calling Captures")
	}
	return it.caps
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
