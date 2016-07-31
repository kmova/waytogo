package main

import (
	"sort"
	"testing"

	"github.com/kmova/waytogo/cli"
)

// Tests if the subcommands of Waytogo are sorted
func TestWaytogoSubcommandsAreSorted(t *testing.T) {
	if !sort.IsSorted(byName(cli.WaytogoCommandUsage)) {
		t.Fatal("WayToGo subcommands are not in sorted order")
	}
}
