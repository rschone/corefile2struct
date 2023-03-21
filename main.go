package main

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
	_ "github.com/rschone/corefile2struct/plugin/my_plugin"
)

var directives = []string{
	"my_plugin",
}

func init() {
	dnsserver.Directives = directives
}

func main() {
	coremain.Run()
}
