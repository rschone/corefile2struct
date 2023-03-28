package myplugin

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type MyPluginCfg struct {
	Arguments     []string
	Core          SubStruct     `cf:"core"`
	Jack          *SubStruct    `cf:"jack"`
	Timeout       time.Duration `cf:"timeout"`
	NumMaxRetries int           `cf:"numMaxRetries" default:"100" check:"nonempty"`
	IPv4          net.IP        `cf:"IPv4" default:"127.0.0.1"`
	Flags         []string      `cf:"flags"`
	Nums2         []int         `cf:"nums"`
	Str           string        `cf:"str" check:"oneof(moje|tvoje|jeho)"`
	Boolean       bool          `cf:"boolean"`
	small         int           `cf:"small"` // cannot be parsed - has to be exported!
}

func (c *MyPluginCfg) Init() error {
	c.Str = "ahoj"
	fmt.Println("cau")
	return nil
}

func (c MyPluginCfg) Check() error {
	if c.Str != "moje" {
		return errors.New("fuck")
	}
	return nil
}

type SubStruct struct {
	KeyPassword string `cf:"key-password"`
	Brokers     string `cf:"brokers"`
	Tls         bool   `cf:"tls"`
}
