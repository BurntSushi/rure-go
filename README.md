Go bindings to RUst's REgex engine
==================================
This package provides cgo bindings to Rust's regex engine.

This is currently a work in progress.


### Documentation

[godoc.org/github.com/BurntSushi/rure-go](http://godoc.org/github.com/BurntSushi/rure-go)


### Benchmarks

These are the benchmarks as they are defined in Go's `regexp` package. Two
important comments:

1. Rust's regex engine hits a performance ceiling because of cgo function call
   overhead. As a result, Go is faster on some benchmarks.
2. Several of the benchmarks are dramatically distorted because of
   optimizations in Rust's regex engine. The benchmarks need to be fine tuned
   to give a more accurate picture. (In fact, none of the benchmarks give an
   accurate depiction of Rust's regex engine overall speed. For example, the
   Hard1 benchmark doesn't even use the regex engine at all. The other
   {easy,medium,hard} benchmarks run the regex engine in reverse from the end
   of the string since the regexes are anchored at the end.)

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
