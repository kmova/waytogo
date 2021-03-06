package cli

import (
	"fmt"

	"github.com/kmova/waytogo/pkg/spf13/cobra"
)

// FlagErrorFunc prints an error messages which matches the format of the
// kmova/waytogo/cli error messages
func FlagErrorFunc(cmd *cobra.Command, err error) error {
	if err == nil {
		return err
	}

	usage := ""
	if cmd.HasSubCommands() {
		usage = "\n\n" + cmd.UsageString()
	}
	return fmt.Errorf("%s\nSee '%s --help'.%s", err, cmd.CommandPath(), usage)
}
