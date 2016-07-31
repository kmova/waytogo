package flags

import flag "github.com/kmova/waytogo/pkg/mflag"

// ClientFlags represents flags for the WayToGo client.
type ClientFlags struct {
	FlagSet   *flag.FlagSet
	Common    *CommonFlags
	PostParse func()

	ConfigDir string
}
