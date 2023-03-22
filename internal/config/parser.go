package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/coredns/caddy"
)

const pluginArgsFieldName = "Arguments"

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

	pointerToStruct := reflect.ValueOf(s)
	structVal := pointerToStruct.Elem()

	if err := p.parsePluginHeader(structVal); err != nil {
		return err
	}

	// if the plugin configuration body is present
	if p.c.Next() {
		if p.c.Val() != "{" {
			return errors.New("'{' expected")
		}
		return p.parseStructure(structVal)
	}
	return nil
}

func (p *Parser) parsePluginHeader(structVal reflect.Value) error {
	pluginName := p.c.Val()
	pluginArgs := p.c.RemainingArgs()
	fmt.Println(pluginName, "args:", pluginArgs)

	if len(pluginArgs) > 0 {
		if err := assignToField(structVal, pluginArgsFieldName, pluginArgs); err != nil {
			return fmt.Errorf("cannot store plugin arguments into field '%s': %w", pluginArgsFieldName, err)
		}
	}
	return nil
}

func (p *Parser) parseStructure(structVal reflect.Value) error {
	log("Parsing a structure..")
	p.applyDefaults(structVal)
	// current symbol in c.Val() is "{"
	for p.c.Next() {
		if p.c.Val() == "}" {
			log("End up the structure!")
			return p.validateStructure(structVal)
		}

		property := p.c.Val()
		propValues := p.c.RemainingArgs()
		if len(propValues) == 0 {
			if p.c.Next() && p.c.Val() == "{" {
				log("structure:", property)
				field := findFieldByTag(structVal, property)
				if field == zeroValue {
					return fmt.Errorf("structure '%s' not found", property)
				}
				if field.Type().Kind() == reflect.Pointer {
					if field.IsNil() {
						newMem := reflect.New(field.Type().Elem())
						field.Set(newMem)
					}
					field = field.Elem()
				}
				p.parseStructure(field)
			} else {
				return fmt.Errorf("property '%s' has no value", property)
			}
		} else {
			field := findFieldByTag(structVal, property)
			if field == zeroValue {
				return errors.New("field not found: " + property)
			}

			value := strings.Join(propValues, ",")
			if err := assignFromString(field, value); err != nil {
				return err
			}
			log("prop:", property, "=", propValues)
		}
	}
	return nil
}

func (p *Parser) applyDefaults(structVal reflect.Value) {
	log("Apply defaults..")
	structType := structVal.Type()
	for i := 0; i < structVal.NumField(); i++ {
		fieldType := structType.Field(i)
		value, ok := fieldType.Tag.Lookup(defaultTag)
		if ok {
			fieldValue := structVal.Field(i)
			assignFromString(fieldValue, value)
		}
	}
}

func (p *Parser) validateStructure(structVal reflect.Value) error {
	log("Validate structure..")
	fmt.Printf("Validating %v\n", structVal)
	// TODO
	return nil
}

func log(a ...any) {
	fmt.Println(a)
}
