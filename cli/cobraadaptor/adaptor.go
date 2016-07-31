package cobraadaptor

import (
	"github.com/kmova/waytogo/api/client"
	"github.com/kmova/waytogo/api/client/system"
	"github.com/kmova/waytogo/cli"
	cliflags "github.com/kmova/waytogo/cli/flags"
	"github.com/kmova/waytogo/pkg/term"
	"github.com/kmova/waytogo/pkg/spf13/cobra"
)

// CobraAdaptor is an adaptor for supporting spf13/cobra commands in the
// waytogo/cli framework
type CobraAdaptor struct {
	rootCmd   *cobra.Command
	waytogoCli *client.WaytogoCli
}

// NewCobraAdaptor returns a new handler
func NewCobraAdaptor(clientFlags *cliflags.ClientFlags) CobraAdaptor {
	stdin, stdout, stderr := term.StdStreams()
	waytogoCli := client.NewWaytogoCli(stdin, stdout, stderr, clientFlags)

	var rootCmd = &cobra.Command{
		Use:           "waytogo [OPTIONS]",
		Short:         "A self-sufficient runtime for containers",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetHelpTemplate(helpTemplate)
	rootCmd.SetFlagErrorFunc(cli.FlagErrorFunc)
	rootCmd.SetOutput(stdout)
	rootCmd.AddCommand(
		system.NewVersionCommand(waytogoCli),
		system.NewInfoCommand(waytogoCli),
	)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	rootCmd.PersistentFlags().MarkShorthandDeprecated("help", "please use --help")

	return CobraAdaptor{
		rootCmd:   rootCmd,
		waytogoCli: waytogoCli,
	}
}

// Usage returns the list of commands and their short usage string for
// all top level cobra commands.
func (c CobraAdaptor) Usage() []cli.Command {
	cmds := []cli.Command{}
	for _, cmd := range c.rootCmd.Commands() {
		if cmd.Name() != "" {
			cmds = append(cmds, cli.Command{Name: cmd.Name(), Description: cmd.Short})
		}
	}
	return cmds
}

func (c CobraAdaptor) run(cmd string, args []string) error {
	if err := c.waytogoCli.Initialize(); err != nil {
		return err
	}
	// Prepend the command name to support normal cobra command delegation
	c.rootCmd.SetArgs(append([]string{cmd}, args...))
	return c.rootCmd.Execute()
}

// Command returns a cli command handler if one exists
func (c CobraAdaptor) Command(name string) func(...string) error {
	for _, cmd := range c.rootCmd.Commands() {
		if cmd.Name() == name {
			return func(args ...string) error {
				return c.run(name, args)
			}
		}
	}
	return nil
}

// GetRootCommand returns the root command. Required to generate the man pages
// and reference docs from a script outside this package.
func (c CobraAdaptor) GetRootCommand() *cobra.Command {
	return c.rootCmd
}

var usageTemplate = `Usage:	{{if not .HasSubCommands}}{{.UseLine}}{{end}}{{if .HasSubCommands}}{{ .CommandPath}} COMMAND{{end}}

{{ .Short | trim }}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{if .HasFlags}}

Options:
{{.Flags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasSubCommands }}

Run '{{.CommandPath}} COMMAND --help' for more information on a command.{{end}}
`

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
