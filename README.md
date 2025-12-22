# gojit

This is a major update to nelhage/gojit - fixing major problems and updating to
modern ABI interactions. It is build for just x86 (amd64) instructions at this time.

The biggest change is proper "jit -> go" func handling. Originally, stack checks
were not handled properly and crashes would occur when gc or stack growth occured.

[Read here for more details.](https://aaronbalke.com/posts/calling-go-functions-from-jit-code/)

This jit works in golang version 1.17+. The proper handling of jit -> go funcs should
continue to work as long as the abi interaction has not major updates.

For handling in golang version 1.16 and earlier please [read](https://www.quasilyte.dev/blog/post/call-go-from-jit/).

Other changes include the removal of bf, and cgo and a simpler build process.
Some intructions have been added based on rasky/gojit.
