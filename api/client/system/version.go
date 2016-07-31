package system

import (
	"runtime"
	"time"

	"golang.org/x/net/context"

	"github.com/kmova/waytogo/api/client"
	"github.com/kmova/waytogo/cli"
	"github.com/kmova/waytogo/version"
	"github.com/kmova/waytogo/utils/templates"
	"github.com/kmova/waytogo/types"
	"github.com/kmova/waytogo/pkg/spf13/cobra"
)

var versionTemplate = `Client:
 Version:      {{.Client.Version}}
 API version:  {{.Client.APIVersion}}
 Go version:   {{.Client.GoVersion}}
 Git commit:   {{.Client.GitCommit}}
 Built:        {{.Client.BuildTime}}
 OS/Arch:      {{.Client.Os}}/{{.Client.Arch}}{{if .Client.Experimental}}
 Experimental: {{.Client.Experimental}}{{end}}{{if .ServerOK}}

Server:
 Version:      {{.Server.Version}}
 API version:  {{.Server.APIVersion}}
 Go version:   {{.Server.GoVersion}}
 Git commit:   {{.Server.GitCommit}}
 Built:        {{.Server.BuildTime}}
 OS/Arch:      {{.Server.Os}}/{{.Server.Arch}}{{if .Server.Experimental}}
 Experimental: {{.Server.Experimental}}{{end}}{{end}}`

type versionOptions struct {
	format string
}

// NewVersionCommand creates a new cobra.Command for `waytogo version`
func NewVersionCommand(waytogoCli *client.WaytogoCli) *cobra.Command {
	var opts versionOptions

	cmd := &cobra.Command{
		Use:   "version [OPTIONS]",
		Short: "Show the WayToGo version information",
		Args:  cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(waytogoCli, &opts)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given go template")

	return cmd
}

func runVersion(waytogoCli *client.WaytogoCli, opts *versionOptions) error {
	ctx := context.Background()

	templateFormat := versionTemplate
	if opts.format != "" {
		templateFormat = opts.format
	}

	tmpl, err := templates.Parse(templateFormat)
	if err != nil {
		return cli.StatusError{StatusCode: 64,
			Status: "Template parsing error: " + err.Error()}
	}

	vd := types.VersionResponse{
		Client: &types.Version{
			Version:      version.Version,
			APIVersion:   waytogoCli.Client().ClientVersion(),
			GoVersion:    runtime.Version(),
			GitCommit:    version.GitCommit,
			BuildTime:    version.BuildTime,
			Os:           runtime.GOOS,
			Arch:         runtime.GOARCH,
		},
	}

	serverVersion, err := waytogoCli.Client().ServerVersion(ctx)
	if err == nil {
		vd.Server = &serverVersion
	}

	// first we need to make BuildTime more human friendly
	t, errTime := time.Parse(time.RFC3339Nano, vd.Client.BuildTime)
	if errTime == nil {
		vd.Client.BuildTime = t.Format(time.ANSIC)
	}

	if vd.ServerOK() {
		t, errTime = time.Parse(time.RFC3339Nano, vd.Server.BuildTime)
		if errTime == nil {
			vd.Server.BuildTime = t.Format(time.ANSIC)
		}
	}

	if err2 := tmpl.Execute(waytogoCli.Out(), vd); err2 != nil && err == nil {
		err = err2
	}
	waytogoCli.Out().Write([]byte{'\n'})
	return err
}
