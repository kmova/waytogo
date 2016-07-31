package system

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/go-units"

	"github.com/kmova/waytogo/api/client"
	"github.com/kmova/waytogo/cli"
	"github.com/kmova/waytogo/pkg/ioutils"
	"github.com/kmova/waytogo/utils"
	"github.com/kmova/waytogo/pkg/spf13/cobra"
)

// NewInfoCommand creates a new cobra.Command for `waytogo info`
func NewInfoCommand(waytogoCli *client.WaytogoCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Display system-wide information",
		Args:  cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(waytogoCli)
		},
	}
	return cmd

}

func runInfo(waytogoCli *client.WaytogoCli) error {
	ctx := context.Background()
	info, err := waytogoCli.Client().Info(ctx)
	if err != nil {
		return err
	}

	ioutils.FprintfIfNotEmpty(waytogoCli.Out(), "Server Version: %s\n", info.ServerVersion)
	ioutils.FprintfIfNotEmpty(waytogoCli.Out(), "Kernel Version: %s\n", info.KernelVersion)
	ioutils.FprintfIfNotEmpty(waytogoCli.Out(), "Operating System: %s\n", info.OperatingSystem)
	ioutils.FprintfIfNotEmpty(waytogoCli.Out(), "OSType: %s\n", info.OSType)
	ioutils.FprintfIfNotEmpty(waytogoCli.Out(), "Architecture: %s\n", info.Architecture)
	fmt.Fprintf(waytogoCli.Out(), "CPUs: %d\n", info.NCPU)
	fmt.Fprintf(waytogoCli.Out(), "Total Memory: %s\n", units.BytesSize(float64(info.MemTotal)))
	ioutils.FprintfIfNotEmpty(waytogoCli.Out(), "ID: %s\n", info.ID)
	fmt.Fprintf(waytogoCli.Out(), "Debug Mode (client): %v\n", utils.IsDebugEnabled())
	fmt.Fprintf(waytogoCli.Out(), "Debug Mode (server): %v\n", info.Debug)


	return nil
}
