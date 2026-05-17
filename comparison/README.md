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

BenchmarkString_Mixed/clipperhouse/displaywidth-8                 6085 ns/op         277.23 MB/s           0 B/op          0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                        9970 ns/op         169.21 MB/s           0 B/op          0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                              19060 ns/op          88.51 MB/s           0 B/op          0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8             6118 ns/op         275.76 MB/s           0 B/op          0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8                   13917 ns/op         121.22 MB/s           0 B/op          0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8                          19263 ns/op          87.58 MB/s           0 B/op          0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8                   54.54 ns/op     2347.10 MB/s           0 B/op          0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                         125.5 ns/op      1020.32 MB/s           0 B/op          0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                               1478 ns/op          86.62 MB/s           0 B/op          0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8                 3265 ns/op         221.74 MB/s           0 B/op          0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8                        5110 ns/op         141.69 MB/s           0 B/op          0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                               7137 ns/op         101.44 MB/s           0 B/op          0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8                   3517 ns/op         479.72 MB/s           0 B/op          0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8                          4746 ns/op         355.48 MB/s           0 B/op          0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8               3454 ns/op         488.36 MB/s           0 B/op          0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8                     11432 ns/op         147.56 MB/s           0 B/op          0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8                    255.5 ns/op       500.88 MB/s           0 B/op          0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8                           264.7 ns/op       483.48 MB/s           0 B/op          0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8                   1320 ns/op         548.44 MB/s           0 B/op          0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8                          2286 ns/op         316.72 MB/s           0 B/op          0 allocs/op

BenchmarkTruncateWithTail/clipperhouse/displaywidth-8             2495 ns/op          70.94 MB/s         192 B/op         14 allocs/op
BenchmarkTruncateWithTail/mattn/go-runewidth-8                    4569 ns/op          38.74 MB/s         192 B/op         14 allocs/op

BenchmarkTruncateWithoutTail/clipperhouse/displaywidth-8          2456 ns/op          93.25 MB/s           0 B/op          0 allocs/op
BenchmarkTruncateWithoutTail/mattn/go-runewidth-8                 5182 ns/op          44.19 MB/s           0 B/op          0 allocs/op
```
