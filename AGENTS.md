The goals and overview of this package can be found in the README.md file,
start by reading that.

The goal of this package is to determine the display (column) width of a
string, UTF-8 bytes, or runes, as would happen in a monospace font, especially
in a terminal.

When troubleshooting, write Go unit tests instead of executing debug scripts.
The tests can return whatever logs or output you need. If those tests are
only for temporary troubleshooting, clean up the tests after the debugging is
done.

(Separate executable debugging scripts are messy, tend to have conflicting
dependencies and are hard to cleanup.)

If you make changes to the trie generation in internal/gen, it can be invoked
by running `go generate` from the top package directory.

We have hard-coded some exceptions to achieve compatibility with go-runewidth.
We consider them technical debt. One example is isExceptionalCombiningMark.
Ideally, we would not have these exceptional cases. Our current theory in the
case of isExceptionalCombiningMark is that go-runewidth is incorrect, but we
don't know for sure.

For PRs, you can use the gh CLI tool to retrieve or post comments.
