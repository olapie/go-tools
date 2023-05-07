package asset

import _ "embed"

//go:embed cert.p12
var CertP12 []byte

//go:embed cert.pem
var CertPEM []byte

//go:embed key.pem
var KeyPEM []byte
