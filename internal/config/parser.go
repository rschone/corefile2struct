package config

import (
	"errors"
	"fmt"
	"k8s.io/utils/strings/slices"
	"reflect"
	"strings"

	"github.com/coredns/caddy"
)

const (
	pluginArgsFieldName = "Arguments"
	cfTag               = "cf"
	defaultTag          = "default"
	checkTag            = "check"
)

type Lexer interface {
	Next() bool
	Val() string
	RemainingArgs() []string
}
type Parser struct {
	lexer Lexer
}

func NewParser(c *caddy.Controller) Parser {
	return Parser{lexer: c}
}

func (p *Parser) Parse(s any) error {
	if !isPointerToStruct(s) {
		return errors.New("invalid argument: pointer to a structure expected")
	}
	if !p.lexer.Next() {
		return errors.New("plugin name expected")
	}

	structVal := reflect.ValueOf(s).Elem()

	if err := p.parsePluginHeader(structVal); err != nil {
		return err
	}

	// if the plugin configuration body is present
	if p.lexer.Next() {
		if p.lexer.Val() != "{" {
			return errors.New("'{' expected")
		}
		return p.parseStructure(structVal)
	}
	return nil
}

func (p *Parser) parsePluginHeader(structVal reflect.Value) error {
	pluginName := p.lexer.Val()
	pluginArgs := p.lexer.RemainingArgs()
	if len(pluginArgs) > 0 {
		if err := assignToField(structVal, pluginArgsFieldName, pluginArgs); err != nil {
			return fmt.Errorf("cannot store plugin '%s' arguments into field '%s': %w", pluginName, pluginArgsFieldName, err)
		}
	}
	return nil
}

func (p *Parser) parseStructure(structVal reflect.Value) error {
	log("Parsing a structure..")
	if err := p.applyDefaults(structVal); err != nil {
		return err
	}

	// current symbol in lexer.Val() is "{"
	for p.lexer.Next() {
		if p.lexer.Val() == "}" {
			log("End up the structure!")
			return p.validateStructure(structVal)
		}

		property := p.lexer.Val()
		propValues := p.lexer.RemainingArgs()
		if len(propValues) == 0 {
			if p.lexer.Next() && p.lexer.Val() == "{" {
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
				if err := p.parseStructure(field); err != nil {
					return err
				}
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

func (p *Parser) applyDefaults(structVal reflect.Value) error {
	log("Apply defaults..")
	structType := structVal.Type()
	for i := 0; i < structVal.NumField(); i++ {
		fieldType := structType.Field(i)
		value, ok := fieldType.Tag.Lookup(defaultTag)
		if ok {
			fieldValue := structVal.Field(i)
			if err := assignFromString(fieldValue, value); err != nil {
				return err
			}
		}
	}
	return executeCustomInit(structVal)
}

func (p *Parser) validateStructure(structVal reflect.Value) error {
	log("Validate structure ", structVal)
	structType := structVal.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag, ok := field.Tag.Lookup(checkTag)
		if ok {
			fieldVal := structVal.Field(i)
			for _, check := range strings.Split(tag, ",") {
				check = strings.Trim(check, " ")
				switch {
				case strings.HasPrefix(check, "nonempty"):
					if err := checkNonEmpty(field.Name, fieldVal); err != nil {
						return err
					}
				case strings.HasPrefix(check, "oneof"):
					start := strings.Index(check, "(")
					end := strings.Index(check, ")")
					content := check[start+1 : end]
					if len(content) > 0 {
						allowedVals := strings.Split(content, "|")
						if !slices.Contains(allowedVals, fieldVal.String()) {
							return fmt.Errorf("value '%s' is not in allowed values %v", fieldVal.String(), allowedVals)
						}
					}
				}
			}
		}
	}

	return executeCustomChecks(structVal)
}

func executeCustomInit(structVal reflect.Value) error {
	if itf, ok := structVal.Addr().Interface().(Initiable); ok && itf != nil {
		if err := itf.Init(); err != nil {
			return err
		}
	}
	return nil
}

type Checkable interface {
	Check() error
}

type Initiable interface {
	Init() error
}

func executeCustomChecks(structVal reflect.Value) error {
	// pointer receiver
	if itf, ok := structVal.Addr().Interface().(Checkable); ok && itf != nil {
		if err := itf.Check(); err != nil {
			return err
		}
	}

	// value receiver
	if itf, ok := structVal.Interface().(Checkable); ok && itf != nil {
		if err := itf.Check(); err != nil {
			return err
		}
	}

	return nil
}

func checkNonEmpty(fieldName string, v reflect.Value) error {
	zero := reflect.Zero(v.Type()).Interface()
	if reflect.DeepEqual(v.Interface(), zero) {
		return fmt.Errorf("non empty check failed for field '%s'", fieldName)
	}
	return nil
}

func log(a ...any) {
	fmt.Println(a)
}

// custom_errors, dns_alias -> nelze (seznamy stejnych properties)
