This is a new Go package for measuring the display (monospace column) width of
strings. It has the same goals as the following packages:

https://github.com/mattn/go-runewidth
https://github.com/jquast/wcwidth

I am developing this package because I believe I can achieve better
performance. Specifically, I intend to operate strictly on strings and []byte,
without ever decoding runes.

The first step is to develop the specification for width measurements. Is there
any official spec? If not, look at the previously mentioned packages above,
and develop a specification for width measurements that is compatible with
them.

After the specification is developed, we will take an approach similar to
https://github.com/clipperhouse/uax29, specifically that we will code-generate
a trie for the table lookups.
