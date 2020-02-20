package rure

// #cgo LDFLAGS: -lrure
//
// #include <stdio.h>
// #include <stdlib.h>
// #include "rure.h"
//
// /*
//  * rure_iter_collect exhausts the given iterator over the given haystack,
//  * and reports all successive non-overlapping match locations in *matches.
//  * *num_matches is set to the number of matches found.
//  *
//  * *matches contains `*num_matches * 2` offsets, where `2 * i` and
//  * `2 * i + 1` represent the start and end byte offsets of match `i`.
//  */
// void rure_iter_collect(rure_iter *it,
//                        const uint8_t *haystack, size_t length,
//                        size_t **matches, size_t *num_matches)
// {
//     rure_match m = {0};
//     size_t len = 0;
//     size_t cap = 64;
//     *matches = malloc(cap * sizeof(size_t));
//     if (NULL == *matches) {
//         fprintf(stderr, "rure_iter_collect: out of memory, aborting\n");
//         abort();
//     }
//     while (rure_iter_next(it, haystack, length, &m)) {
//         if ((len * 2 + 1) >= cap) {
//             cap *= 2;
//             *matches = realloc(*matches, cap * sizeof(size_t));
//             if (NULL == *matches) {
//                 fprintf(
//                     stderr, "rure_iter_collect: out of memory, aborting\n");
//                 abort();
//             }
//         }
//         (*matches)[len * 2 + 0] = m.start;
//         (*matches)[len * 2 + 1] = m.end;
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

// Flags for modifying regex behavior. All flags can be modified in the
// expression itself using standard syntax. e.g., `(?i)` enables case
// insensitivity and `(?-i)` disables it.
const (
	// FlagCaseI is the case insensitive (i) flag.
	FlagCaseI = C.RURE_FLAG_CASEI
	// FlagMulti is the multi-line matching (m) flag.
	// (^ and $ match new line boundaries.)
	FlagMulti = C.RURE_FLAG_MULTI
	// FlagDotNL is the any character (s) flag. (. matches new line.)
	FlagDotNL = C.RURE_FLAG_DOTNL
	// FlagSwapGreed is the greedy swap (U) flag.
	// (e.g., + is ungreedy and +? is greedy.)
	FlagSwapGreed = C.RURE_FLAG_SWAP_GREED
	// FlagSpace is the ignore whitespace (x) flag.
	FlagSpace = C.RURE_FLAG_SPACE
	// FlagUnicode is the Unicode (u) flag.
	FlagUnicode = C.RURE_FLAG_UNICODE
	// FlagDefault is used when calling MustCompile or Compile.
	FlagDefault = FlagUnicode
)

// Regex is a compiled regular expression.
//
// It can be used safely from multiple goroutines simultaneously.
type Regex struct {
	pattern string
	p       *C.rure
}

// Options represents non-flag compile time options.
//
// For example, calling SetSizeLimit will place an upper bound on how big
// the compiled regular expression can be.
type Options struct {
	p *C.rure_options
}

// Captures represents start and end locations for every matching capture group
// in a regular expression match.
//
// It is not safe to use from multiple goroutines simultaneously.
type Captures struct {
	re *Regex
	p  *C.rure_captures
	ok bool
}

// Iter is an iterator over successive non-overlapping matches in a haystack.
//
// It is not safe to use from multiple goroutines simultaneously.
type Iter struct {
	re       *Regex
	p        *C.rure_iter
	haystack []byte
	match    C.rure_match
}

// Error is an error that caused compilation of a regular expression to fail.
//
// Most errors are syntax errors, but an error can be returned if the compiled
// regular expression would be too big.
type Error struct {
	p *C.rure_error
}

// MustCompile is like Compile, but if there was a problem compiling the
// pattern, then it will panic.
func MustCompile(pattern string) *Regex {
	re, err := Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("regex.MustCompile failed on %s: %s", pattern, err))
	}
	return re
}

// Compile is like CompileOptions, but uses default flags (Unicode enabled)
// and default size limits.
//
// If there was a problem compiling the pattern, then no Regex is returned and
// a non-nil error is returned.
func Compile(pattern string) (*Regex, error) {
	return CompileOptions(pattern, FlagDefault, nil)
}

// CompileOptions compiles a pattern (in UTF-8) to a regular expression
// suitable for searching text.
//
// Flags is a bitfield of the Flag constants in this package. A value of `0`
// disables all flags.
//
// Options is a set non-flag configuration settings for the compiled regular
// expression. When set to nil, default settings are used.
//
// If there was a problem compiling the pattern (including if it is not valid
// UTF-8), then an error is returned.
//
// N.B. When disabling Unicode support (either via an explicit flag here or in
// the pattern), one can use escape sequences to match bytes that aren't valid
// UTF-8. For example, use `\\xFF` instead of `\xFF`.
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
		asUint8Ptr(noCopyBytes(pattern)),
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

// IsMatch returns true if text matches re.
func (re *Regex) IsMatch(text string) bool {
	return re.IsMatchBytesAt(noCopyBytes(text), 0)
}

// IsMatchBytes returns true if text matches re.
func (re *Regex) IsMatchBytes(text []byte) bool {
	return re.IsMatchBytesAt(text, 0)
}

// IsMatchAt returns true if text matches re starting at index i.
func (re *Regex) IsMatchAt(text string, i int) bool {
	return re.IsMatchBytesAt(noCopyBytes(text), i)
}

// IsMatchBytesAt returns true if text matches re starting at index i.
func (re *Regex) IsMatchBytesAt(text []byte, i int) bool {
	return bool(C.rure_is_match(
		re.p, asUint8Ptr(text), C.size_t(len(text)), C.size_t(i)))
}

// ShortestMatch returns the end location of a match in text if it exists. This
// may return an end location that occurs before the end of the proper
// leftmost-first match.
//
// If no match exists, false is returned.
//
// For example, matching `a+` against `aaaaa` will return `1` (while `Find`
// will report `5` as the end).
func (re *Regex) ShortestMatch(text string) (end int, ok bool) {
	return re.ShortestMatchBytes(noCopyBytes(text))
}

// ShortestMatchBytes returns the end location of a match in text if it exists.
// This may return an end location that occurs before the end of the proper
// leftmost-first match.
//
// If no match exists, false is returned.
//
// For example, matching `a+` against `aaaaa` will return `1` (while `Find`
// will report `5` as the end).
func (re *Regex) ShortestMatchBytes(text []byte) (end int, ok bool) {
	var cend C.size_t
	ok = bool(C.rure_shortest_match(
		re.p, asUint8Ptr(text), C.size_t(len(text)), 0, &cend))
	end = int(cend)
	return
}

// Find returns the start and end location of the leftmost-first match in text
// if it exists.
//
// If no match exists, false is returned.
func (re *Regex) Find(text string) (start, end int, ok bool) {
	return re.FindBytes(noCopyBytes(text))
}

// FindBytes returns the start and end location of the leftmost-first match in
// text if it exists.
//
// If no match exists, false is returned.
func (re *Regex) FindBytes(text []byte) (start, end int, ok bool) {
	match := C.rure_match{}
	ok = bool(C.rure_find(
		re.p, asUint8Ptr(text), C.size_t(len(text)), 0, &match))
	if ok {
		start, end = int(match.start), int(match.end)
	}
	return
}

// FindAll returns all successive non-overlapping matches of re in text.
//
// The slice returned contains a pair of start and end offsets for each match
// found. The start and end offset for match i is indexed by i*2 and i*2+1,
// respectively.
//
// This may be faster than using Iter since the slice of matches is built in
// C code.
func (re *Regex) FindAll(text string) []int {
	return re.FindAllBytes(noCopyBytes(text))
}

// FindAllBytes returns all successive non-overlapping matches of re in text.
//
// The slice returned contains a pair of start and end offsets for each match
// found. The start and end offset for match i is indexed by i*2 and i*2+1,
// respectively.
//
// This may be faster than using Iter since the slice of matches is built in
// C code.
func (re *Regex) FindAllBytes(text []byte) []int {
	it := C.rure_iter_new(re.p)
	defer C.rure_iter_free(it)

	haystack := asUint8Ptr(text)
	len := C.size_t(len(text))
	nmatches := C.size_t(0)
	matches := (*C.size_t)(nil)
	defer func() {
		if matches != nil {
			C.free(unsafe.Pointer(matches))
		}
	}()

	C.rure_iter_collect(it, haystack, len, &matches, &nmatches)
	if nmatches == 0 {
		return nil
	}

	// Copy the matches from C memory to Go memory.
	matchesInts := make([]int, 2*nmatches)
	p := uintptr(unsafe.Pointer(matches))
	stride := unsafe.Sizeof(C.size_t(0))
	for i := uintptr(0); i < uintptr(nmatches); i++ {
		base := i * 2
		start := unsafe.Pointer(p + (base * stride))
		end := unsafe.Pointer(p + ((base + 1) * stride))
		matchesInts[base] = int(*(*C.size_t)(start))
		matchesInts[base+1] = int(*(*C.size_t)(end))
	}
	return matchesInts
}

// NewCaptures allocates room for storing the start and end offset of each
// capturing group in re.
//
// Captures may be reused in subsequent calls. When it is reused, its internal
// state is reset.
//
// Captures may not be used from multiple threads simultaneously.
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

// Captures populates caps with the start and end locations of all matching
// capturing groups in re for text.
//
// If no match is found, then false is returned.
//
// A caps value can be built with re.NewCaptures().
//
// caps must not be nil.
func (re *Regex) Captures(caps *Captures, text string) bool {
	return re.CapturesBytes(caps, noCopyBytes(text))
}

// CapturesBytes populates caps with the start and end locations of all
// matching capturing groups in re for text.
//
// If no match is found, then false is returned.
//
// A caps value can be built with re.NewCaptures().
//
// caps must not be nil.
func (re *Regex) CapturesBytes(caps *Captures, text []byte) bool {
	caps.ok = bool(C.rure_find_captures(
		re.p, asUint8Ptr(text), C.size_t(len(text)), 0, caps.p))
	return caps.ok
}

// Iter returns an iterator over successive non-overlapping matches of re
// in text.
//
// Next must be called on the iterator before accessing match information.
func (re *Regex) Iter(text string) *Iter {
	return re.IterBytes(noCopyBytes(text))
}

// IterBytes returns an iterator over successive non-overlapping matches of re
// in text.
//
// Next must be called on the iterator before accessing match information.
func (re *Regex) IterBytes(text []byte) *Iter {
	return newIter(re, text)
}

// CaptureNames returns a slice of the names of call capturing groups in this
// regex. The slice has the same order as the order of the appearance of each
// capturing group. Index 0 corresponds to the entire regex match, and is
// therefore always unnamed. Unnamed capturing groups are always represented by
// an empty string.
func (re *Regex) CaptureNames() []string {
	it := C.rure_iter_capture_names_new(re.p)
	defer C.rure_iter_capture_names_free(it)

	var names []string
	var name *C.char
	for C.rure_iter_capture_names_next(it, &name) {
		names = append(names, C.GoString(name))
	}
	return names
}

// NewOptions returns a fresh options value for configuring non-flag options
// of a regex.
//
// Options can be passed to multiple calls to CompileOptions.
//
// Options is not safe to mutate from multiple goroutines simultaneously, but
// it may be used in calls to CompileOptions from multiple goroutines
// simultaneously.
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

// SetSizeLimit sets the approximate size limit (in bytes) of the compiled
// regular expression.
//
// If a pattern would result in a compiled program large than this size, then
// compilation will return an error.
func (opts *Options) SetSizeLimit(limit int) {
	C.rure_options_size_limit(opts.p, C.size_t(limit))
}

// SetDFASizeLimit sets the approximate size limit (in bytes) of the DFA's
// cache size. It is a per-thread limit, so that a Regex used in multiple
// threads simultaneously may use up to this many bytes per-thread of usage.
//
// 0 is a legal value.
func (opts *Options) SetDFASizeLimit(limit int) {
	C.rure_options_dfa_size_limit(opts.p, C.size_t(limit))
}

// IsMatch returns true if caps corresponds to a match in a regular expression.
func (caps *Captures) IsMatch() bool {
	return caps.ok
}

// Group returns the start and end offsets for the capturing group indexed by
// i. Capturing groups are indexed by the appearance of their opening
// parenthesis in the pattern.
//
// If the capturing group was not part of the match, then this returns false.
//
// Note that capture group 0 always corresponds to the full match of the
// regular expression and is always unnamed.
func (caps *Captures) Group(i int) (start, end int, ok bool) {
	match := C.rure_match{}
	ok = bool(C.rure_captures_at(caps.p, C.size_t(i), &match))
	if ok {
		start, end = int(match.start), int(match.end)
	}
	return
}

// GroupName is like Group, but uses the name of a capturing group instead of
// its index. Named capture groups look like (?P<foo>re) in the pattern.
//
// If no such named capture group exists or if it wasn't part of the match
// of the regular expression, GroupName returns false.
func (caps *Captures) GroupName(name string) (start, end int, ok bool) {
	i := C.rure_capture_name_index(caps.re.p, C.CString(name))
	if i == -1 {
		return
	}
	return caps.Group(int(i))
}

// Len returns the number of capturing groups.
//
// Once caps is created, this never changes.
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

// Next advances the iterator. If it finds a match, it returns true, and
// otherwise returns false. Once it returns false, it will always return false.
//
// This must be called before a call to Match.
//
// If caps is nil, then capture information is not extracted. If caps is not
// nil, then the start and end offsets of each matching capturing group are
// stored in caps.
func (it *Iter) Next(caps *Captures) bool {
	haystack := asUint8Ptr(it.haystack)
	len := C.size_t(len(it.haystack))

	if caps == nil {
		return bool(C.rure_iter_next(it.p, haystack, len, &it.match))
	}
	caps.ok = bool(C.rure_iter_next_captures(it.p, haystack, len, caps.p))
	C.rure_captures_at(caps.p, 0, &it.match)
	return caps.ok
}

// Match returns the start and end offsets of the current match in the
// iteraotr.
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

// Converts a byte slice to a *C.uint8_t.
//
// This works even for empty slices.
func asUint8Ptr(bs []byte) *C.uint8_t {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	return (*C.uint8_t)(unsafe.Pointer(bh.Data))
}
