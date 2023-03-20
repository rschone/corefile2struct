package config

import (
	"fmt"
	"github.com/coredns/caddy"
)

type Parser struct {
}

func (p *Parser) Parse(c *caddy.Controller, v any) error {

	for c.Next() {
		fmt.Println(c.Val())
	}

	return nil
}
