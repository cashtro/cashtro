package server

import _ "embed"

//go:embed index.html
var indexHTML []byte

//go:embed favicon.svg
var faviconSVG []byte
