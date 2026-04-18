## Compatibility

In real-world text, you should see the same outputs from
`clipperhouse/displaywidth`, `mattn/go-runewidth`, and `rivo/uniseg`.

The tests in this `comparison` package exercise the behaviors of the three
libraries. Extensive details are available in the
[compatibility analysis](COMPATIBILITY_ANALYSIS.md).

## Benchmarks

```bash
go test -bench=. -benchmem
```

```
goos: darwin
goarch: arm64
pkg: github.com/clipperhouse/displaywidth/comparison
cpu: Apple M2

BenchmarkString_Mixed/clipperhouse/displaywidth-8              6326 ns/op         266.66 MB/s           0 B/op          0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                     9984 ns/op         168.97 MB/s           0 B/op          0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                           19602 ns/op          86.06 MB/s           0 B/op          0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8          6167 ns/op         273.55 MB/s           0 B/op          0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8                14022 ns/op         120.31 MB/s           0 B/op          0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8                       19608 ns/op          86.04 MB/s           0 B/op          0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8                60.62 ns/op     2111.60 MB/s           0 B/op          0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                      122.3 ns/op      1047.01 MB/s           0 B/op          0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                            1490 ns/op          85.89 MB/s           0 B/op          0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8              3313 ns/op         218.51 MB/s           0 B/op          0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8                     5009 ns/op         144.55 MB/s           0 B/op          0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                            6868 ns/op         105.42 MB/s           0 B/op          0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8                3430 ns/op         491.90 MB/s           0 B/op          0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8                       4833 ns/op         349.09 MB/s           0 B/op          0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8            3494 ns/op         482.77 MB/s           0 B/op          0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8                  11724 ns/op         143.89 MB/s           0 B/op          0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8                 256.0 ns/op       500.02 MB/s           0 B/op          0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8                        265.0 ns/op       483.01 MB/s           0 B/op          0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8                1381 ns/op         524.30 MB/s           0 B/op          0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8                       2345 ns/op         308.70 MB/s           0 B/op          0 allocs/op

BenchmarkTruncateWithTail/clipperhouse/displaywidth-8          2755 ns/op          64.24 MB/s         192 B/op         14 allocs/op
BenchmarkTruncateWithTail/mattn/go-runewidth-8                 4683 ns/op          37.80 MB/s         192 B/op         14 allocs/op

BenchmarkTruncateWithoutTail/clipperhouse/displaywidth-8       2481 ns/op          92.30 MB/s           0 B/op          0 allocs/op
BenchmarkTruncateWithoutTail/mattn/go-runewidth-8              5334 ns/op          42.93 MB/s           0 B/op          0 allocs/op
```
