module github.com/ije/esbuild-internal

// Support for Go 1.13 is deliberate so people can build esbuild
// themselves for old OS versions. Please do not change this.
go 1.24.0

require (
	github.com/k0kubun/pp v3.0.1+incompatible
	// This dependency cannot be upgraded or esbuild would no longer
	// compile with Go 1.13. Please do not change this. For more info,
	// please read this: https://esbuild.github.io/faq/#old-go-version
	golang.org/x/sys v0.38.0
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
)
