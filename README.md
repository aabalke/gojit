# aabalke/gojit: Jit Compiler for Golang

aabalke/gojit is a revival of nelhage/gojit and rasky/gojit, providing jit compilation for Golang version 1.17+. The major functionality this fork adds is the ability for Go functions to be called from jit code, without causing errors from stack checks during garbage collection or stack growth.

# Resources
[Writeup](https://aaronbalke.com/posts/gojit/)

[Video](https://youtu.be/7fGyjb-6mYA)

# Requirements

This jit compiler requires an x86/amd64 system, and go version 1.17+.
For handling in golang version 1.16 and earlier please [read](https://www.quasilyte.dev/blog/post/call-go-from-jit/).

If a Go function is called from jit code, the Go functions must have a nosplit directive.
If the JIT mutates variables though pointers, those variables must be heap allocated, and/or global.

# Examples

See the example directory for basic examples, or see the [Guac Emulator](https://github.com/aabalke/guac) if you are interested in a real-life example. It includes the setup of multiple JIT pages, analysis and metrics for JIT compilation thresholds, invalidation, and an LRU cache for invalidating dead JIT compiler blocks.
