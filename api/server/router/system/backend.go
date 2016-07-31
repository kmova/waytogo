package system

import (
	"github.com/kmova/waytogo/types"
)

// Backend is the methods that need to be implemented to provide
// system specific functionality.
type Backend interface {
	SystemInfo() (*types.Info, error)
	SystemVersion() types.Version
}
