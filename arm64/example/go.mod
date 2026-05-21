module example

replace github.com/aabalke/gojit/arm64 => ../../arm64

go 1.26.2

require github.com/aabalke/gojit/arm64 v0.0.0-00010101000000-000000000000

require (
	github.com/edsrzf/mmap-go v1.2.0 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
)
