package config

import (
	"errors"
	"fmt"

	"github.com/coredns/caddy"
)

type Parser struct {
	c *caddy.Controller
}

func NewParser(c *caddy.Controller) Parser {
	return Parser{c: c}
}

func (p *Parser) Parse(s any) error {
	if !isPointerToStruct(s) {
		return errors.New("invalid argument: pointer to a structure expected")
	}
	if !p.c.Next() {
		return errors.New("plugin name expected")
	}

	if err := p.parsePluginHeader(s); err != nil {
		return err
	}

	// if the plugin configuration body is present
	if p.c.Next() {
		if p.c.Val() != "{" {
			return errors.New("'{' expected")
		}
		return p.parseStructure(s)
	}
	return nil
}

func (p *Parser) parsePluginHeader(s any) error {
	pluginName := p.c.Val()
	pluginArgs := p.c.RemainingArgs()
	fmt.Println(pluginName, "args:", pluginArgs)
	// TODO: place args into v.Arguments, if both present
	return nil
}

func (p *Parser) parseStructure(s any) error {
	fmt.Println("Parsing a structure..")
	p.applyDefaults(s)
	// current symbol in c.Val() is "{"
	for p.c.Next() {
		if p.c.Val() == "}" {
			fmt.Println("End up the structure!")
			return p.validateStructure(s)
		}

		propName := p.c.Val()
		propValues := p.c.RemainingArgs()
		if len(propValues) == 0 {
			if p.c.Next() && p.c.Val() == "{" {
				fmt.Println("structure:", propName)
				// TODO: find field with propName
				p.parseStructure(s)
			} else {
				// TODO: property with no value?
			}
		} else {
			fmt.Println("prop:", propName, "=", propValues)
			// TODO: assign value of the property to an appropriate field in the structure
		}
	}
	return nil
}

func (p *Parser) applyDefaults(s any) {
	fmt.Println("Apply defaults..")
	// TODO
}

func (p *Parser) validateStructure(s any) error {
	fmt.Println("Validate structure..")
	// TODO
	return nil
}
