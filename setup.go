package neoteqts4via6

import (
	"context"
	"net"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type NeoteqTS4via6 struct {
	Next        plugin.Handler
	Fallthrough bool
}

func init() {
	plugin.Register("neoteqts4via6", setup)
}

func setup(c *caddy.Controller) error {
	p := NeoteqTS4via6{}
	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "fallthrough":
				p.Fallthrough = true
			default:
				return c.Errf("unknown property '%s'", c.Val())
			}
		}
	}
	return nil
}

func (p NeoteqTS4via6) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	qName := r.Question[0].Name
	qType := r.Question[0].Qtype

	if qType == dns.TypeAAAA {
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
	}

	if qType == dns.TypeA {
		msg := new(dns.Msg)
		msg.SetReply(r)
		msg.Authoritative = true
		msg.Rcode = dns.RcodeSuccess // Explizit NOERROR setzen
		w.WriteMsg(msg)
		return dns.RcodeSuccess, nil
	}

	if p.Fallthrough {
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}

	return dns.RcodeNameError, nil
}

func (p NeoteqTS4via6) Name() string {
	return "neoteqts4via6"
}
