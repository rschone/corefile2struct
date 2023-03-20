package myplugin

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type handler struct {
	Next plugin.Handler
}

func (h *handler) ServeDNS(ctx context.Context, writer dns.ResponseWriter, requestMsg *dns.Msg) (int, error) {
	return plugin.NextOrFailure(h.Name(), h.Next, ctx, writer, requestMsg)
}

func (h *handler) Name() string {
	return pluginName
}
