// +build !daemon

package main

import (
	"strings"
	"testing"
)

func TestCmdDaemon(t *testing.T) {
	proxy := NewDaemonProxy()
	err := proxy.CmdDaemon("--help")
	if err == nil {
		t.Fatal("Expected CmdDaemon to fail on Windows.")
	}

	if !strings.Contains(err.Error(), "Please run `waytogod`") {
		t.Fatalf("Expected an error about running waytogod, got %s", err)
	}
}
