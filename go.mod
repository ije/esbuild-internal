module github.com/ije/esbuild-internal

// Support for Go 1.13 is deliberate so people can build esbuild
// themselves for old OS versions. Please do not change this.
go 1.24.0

// This dependency cannot be upgraded or esbuild would no longer
// compile with Go 1.13. Please do not change this. For more info,
// please read this: https://esbuild.github.io/faq/#old-go-version
require golang.org/x/sys v0.38.0
