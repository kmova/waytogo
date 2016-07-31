package cli

// Command is the struct containing the command name and description
type Command struct {
	Name        string
	Description string
}

// WaytogoCommandUsage lists the top level waytogo commands and their short usage
var WaytogoCommandUsage = []Command{
	{"inspect", "Return low-level information on a directory"},
}

// WaytogoCommands stores all the waytogo command
var WaytogoCommands = make(map[string]Command)

func init() {
	for _, cmd := range WaytogoCommandUsage {
		WaytogoCommands[cmd.Name] = cmd
	}
}
