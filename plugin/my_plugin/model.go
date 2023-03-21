package myplugin

import (
	"net"
	"time"
)

type MyPluginCfg struct {
	Core          SubStruct     `cf:"core"`
	Jack          *SubStruct    `cf:"jack"`
	Timeout       time.Duration `cf:"timeout"`
	NumMaxRetries int           `cf:"numMaxRetries,lt(100),{100}"`
	IPv4          net.IP        `cf:"IPv4,{127.0.0.1}"`
	Flags         []string      `cf:"flags"`
	Str           string        `cf:"str"`
	Boolean       bool          `cf:"boolean"`
	small         int           `cf:"small"`
}

type SubStruct struct {
	KeyPassword string `cf:"key-password,{mypassword}"`
	Brokers     string `cf:"brokers"`
	Tls         bool   `cf:"tls"`
}
