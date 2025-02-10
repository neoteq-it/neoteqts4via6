package neoteqts4via6

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type NeoteqTS4via6 struct {
	Next plugin.Handler
}

func init() {
	plugin.Register("neoteqts4via6", setup)
}

func setup(c *caddy.Controller) error {
	p := &NeoteqTS4via6{}

	c.Next() // Bewegt den Parser weiter (Pflicht in CoreDNS)

	c.OnStartup(func() error {
		fmt.Println("NeoteqTS4via6 Plugin geladen!")
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
		// Keine Antwort setzen, was zu "No Answer" führt
		w.WriteMsg(msg)
		return dns.RcodeSuccess, nil
	}

	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}

func (p NeoteqTS4via6) Name() string {
	return "neoteqts4via6"
}
