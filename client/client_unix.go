// +build linux freebsd solaris openbsd

package client

// DefaultWaytogoHost defines os specific default if WAYTOGO_HOST is unset
const DefaultWaytogoHost = "unix:///var/run/waytogo.sock"
