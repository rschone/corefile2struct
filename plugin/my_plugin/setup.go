package myplugin

import (
	"fmt"
	"github.com/coredns/caddy"
	config2 "github.com/rschone/corefile2struct/internal/config_parser"
)

const pluginName = "my_plugin"

func init() {
	caddy.RegisterPlugin(pluginName, caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	var parser config2.Parser
	var cfg MyPluginCfg
	err := parser.Parse(c, &cfg)
	if err != nil {
		return err
	}

	fmt.Println(cfg)

	return nil
}
