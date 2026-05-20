// Package baseasset embeds the shared base.html template for use by
// command-line entrypoints.
package baseasset

import _ "embed"

//go:embed base.html
var HTML []byte
