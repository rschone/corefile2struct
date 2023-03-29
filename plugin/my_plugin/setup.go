package myplugin

import (
	"fmt"
	"github.com/rschone/corefile2struct/internal/corefile"

	"github.com/coredns/caddy"
)

const pluginName = "my_plugin"

func init() {
	caddy.RegisterPlugin(pluginName, caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	var myCfg MyPluginCfg
	if err := corefile.Parse(c, &myCfg); err != nil {
		return err
	}

	fmt.Println("Parsed cfg:", myCfg)

	return nil
}
