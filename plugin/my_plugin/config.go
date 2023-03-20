package myplugin

import (
	"net"
	"time"

	"github.com/coredns/caddy"
)

type config struct {
	Core          Block
	Jack          *Block
	Timeout       time.Duration
	NumMaxRetries int
	IPv4          net.IP
}

type Block struct {
	KeyPassword string
	Brokers     string
	Tls         bool
}

func parse(c *caddy.Controller) (config, error) {
	var cacheParsed, blockParsed bool
	cfg := config{}
	//var err error

	if c.Next() {
		args := c.RemainingArgs()
		if len(args) > 0 {
			return cfg, c.Errf("unknown plugin parameters '%v'", args)
		}

		//for c.NextBlock() {
		//	switch c.Val() {
		//	case "aws":
		//		if cfg.aws, err = parseAws(c); err != nil {
		//			return cfg, err
		//		}
		//	case "device-context-cache":
		//		if cfg.cacheConfig, err = cache.ParseCacheConfig(c, log, "device-context"); err != nil {
		//			return cfg, err
		//		}
		//		cacheParsed = true
		//	case "block-response":
		//		rcfg, err := block.ParseResponseConfig(c, log)
		//		if err != nil {
		//			return cfg, err
		//		}
		//		cfg.blockResponse = &rcfg
		//		blockParsed = true
		//	default:
		//		log.Warningf("Unknown property '%s'.", c.Val())
		//	}
		//}

		if !cacheParsed {
			return cfg, c.Err("cache configuration is required")
		}
		if !blockParsed {
			return cfg, c.Err("block  configuration is required")
		}
	}

	return cfg, nil
}

//
//func parseAws(c *caddy.Controller) (awsConfig, error) {
//	cfg := awsConfig{}
//	var err error
//
//	for c.Next() {
//		switch c.Val() {
//		case "configFile":
//			c.Next()
//			cfg.configFile = c.Val()
//			log.Infof("aws.configFile: %s.", cfg.configFile)
//		case "credentialsFile":
//			c.Next()
//			cfg.credentialsFile = c.Val()
//			log.Infof("aws.credentialsFile: %s.", cfg.credentialsFile)
//		case "endpoint":
//			if c.NextArg() {
//				cfg.endpoint = c.Val()
//				log.Infof("aws.endpoint: %s.", cfg.endpoint)
//			}
//		case "retryer":
//			if cfg.retryer, err = parseRetryer(c); err != nil {
//				return cfg, err
//			}
//		case "timeout":
//			c.Next()
//			if cfg.timeout, err = time.ParseDuration(c.Val()); err != nil {
//				return cfg, err
//			}
//			log.Infof("aws.timeout: %v.", cfg.timeout)
//		case "{":
//		case "}":
//			return cfg, nil
//		default:
//			log.Warningf("Unknown property '%s'.", c.Val())
//		}
//	}
//
//	return cfg, errors.New("wrong configuration parameter lines")
//}
//
//func parseRetryer(c *caddy.Controller) (client.DefaultRetryer, error) {
//	cfg := client.DefaultRetryer{}
//	var err error
//
//	for c.Next() {
//		switch c.Val() {
//		case "numMaxRetries":
//			c.Next()
//			if cfg.NumMaxRetries, err = strconv.Atoi(c.Val()); err != nil {
//				return cfg, err
//			}
//			log.Infof("aws.retryer.numMaxRetries: %v.", cfg.NumMaxRetries)
//		case "minRetryDelay":
//			c.Next()
//			if cfg.MinRetryDelay, err = time.ParseDuration(c.Val()); err != nil {
//				return cfg, err
//			}
//			log.Infof("aws.retryer.minRetryDelay: %d.", cfg.MinRetryDelay)
//		case "minThrottleDelay":
//			c.Next()
//			if cfg.MinThrottleDelay, err = time.ParseDuration(c.Val()); err != nil {
//				return cfg, err
//			}
//			log.Infof("aws.retryer.minThrottleDelay: %v.", cfg.MinThrottleDelay)
//		case "maxRetryDelay":
//			c.Next()
//			if cfg.MaxRetryDelay, err = time.ParseDuration(c.Val()); err != nil {
//				return cfg, err
//			}
//			log.Infof("aws.retryer.maxRetryDelay: %v.", cfg.MaxRetryDelay)
//		case "maxThrottleDelay":
//			c.Next()
//			if cfg.MaxThrottleDelay, err = time.ParseDuration(c.Val()); err != nil {
//				return cfg, err
//			}
//			log.Infof("aws.retryer.maxThrottleDelay: %v.", cfg.MaxThrottleDelay)
//		case "{":
//		case "}":
//			return cfg, nil
//		default:
//			log.Warningf("Unknown property '%s'.", c.Val())
//		}
//	}
//
//	return cfg, errors.New("wrong configuration parameter lines")
//}
