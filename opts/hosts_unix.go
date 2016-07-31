// +build !windows

package opts

import "fmt"

// DefaultHost constant defines the default host string used by waytogo on other hosts than Windows
var DefaultHost = fmt.Sprintf("unix://%s", DefaultUnixSocket)
