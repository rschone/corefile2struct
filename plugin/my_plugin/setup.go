package myplugin

import (
	"fmt"

	"github.com/coredns/caddy"
	cfg "github.com/rschone/corefile2struct/internal/config"
)

const pluginName = "my_plugin"

func init() {
	caddy.RegisterPlugin(pluginName, caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	parser := cfg.NewParser(c)

	var myCfg MyPluginCfg
	if err := parser.Parse(&myCfg); err != nil {
		return err
	}

	fmt.Println("Parsed cfg:", myCfg)

	return nil
}
