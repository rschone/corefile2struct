package myplugin

import (
	"net"
	"time"
)

type MyPluginCfg struct {
	Arguments     []string
	Core          SubStruct     `cf:"core"`
	Jack          *SubStruct    `cf:"jack"`
	Timeout       time.Duration `cf:"timeout"`
	NumMaxRetries int           `cf:"numMaxRetries" default:"100" cond:"notempty,nz,lt(h),gt(12)"`
	IPv4          net.IP        `cf:"IPv4" default:"127.0.0.1"`
	Flags         []string      `cf:"flags"`
	Nums2         []int         `cf:"nums"`
	Str           string        `cf:"str"`
	Boolean       bool          `cf:"boolean"`
	small         int           `cf:"small"`
}

type SubStruct struct {
	KeyPassword string `cf:"key-password"`
	Brokers     string `cf:"brokers"`
	Tls         bool   `cf:"tls"`
}
