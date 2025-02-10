package neoteqts4via6

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

type NeoteqTS4via6 struct {
	Next        plugin.Handler
	Fallthrough bool
	TTL         uint32
}

func init() {
	plugin.Register("neoteqts4via6", setup)
}

func setup(c *caddy.Controller) error {
	p := &NeoteqTS4via6{
		TTL: 60, // Default TTL
	}

	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "fallthrough":
				p.Fallthrough = true
			case "ttl":
				if !c.NextArg() {
					return plugin.Error("neoteqts4via6", c.ArgErr())
				}
				ttl, err := strconv.ParseUint(c.Val(), 10, 32)
				if err != nil {
					return plugin.Error("neoteqts4via6", c.Errf("Invalid TTL: %s", c.Val()))
				}
				p.TTL = uint32(ttl)
			default:
				return plugin.Error("neoteqts4via6", c.Errf("Unknown property '%s'", c.Val()))
			}
		}
	}

	c.OnStartup(func() error {
		log.Info("NeoteqTS4via6 Plugin loaded")
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		p.Next = next
		return p
	})

	return nil
}

func (p NeoteqTS4via6) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	qName := r.Question[0].Name
	qType := r.Question[0].Qtype

	switch qType {
	case dns.TypeAAAA:
		ipv6, err := ResolveIPv6(strings.TrimSuffix(qName, "."))
		if err == nil {
			msg := new(dns.Msg)
			msg.SetReply(r)
			msg.Authoritative = true

			rr := &dns.AAAA{
				Hdr:  dns.RR_Header{Name: qName, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 60},
				AAAA: net.ParseIP(ipv6),
			}
			msg.Answer = append(msg.Answer, rr)
			w.WriteMsg(msg)
			return dns.RcodeSuccess, nil
		}
		// Fallthrough, wenn keine IPv6-Adresse gefunden wurde
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)

	case dns.TypeA:
		// Leere Antwort für A-Anfragen, dann Fallthrough
		msg := new(dns.Msg)
		msg.SetReply(r)
		msg.Authoritative = true
		msg.Rcode = dns.RcodeSuccess
		w.WriteMsg(msg)
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)

	default:
		// Für alle anderen Anfragen direkt Fallthrough
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}
}

func (p *NeoteqTS4via6) Name() string {
	return "neoteqts4via6"
}
