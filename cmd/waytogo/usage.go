package main

import (
	"sort"

	"github.com/kmova/waytogo/cli"
)

type byName []cli.Command

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func sortCommands(commands []cli.Command) []cli.Command {
	waytogoCommands := make([]cli.Command, len(commands))
	copy(waytogoCommands, commands)
	sort.Sort(byName(waytogoCommands))
	return waytogoCommands
}
