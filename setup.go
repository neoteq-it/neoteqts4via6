package neoteqts4via6

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// NeoteqTS4via6 ist das CoreDNS-Plugin
type NeoteqTS4via6 struct {
	Next plugin.Handler
}

// init registriert das Plugin bei CoreDNS
func init() {
	plugin.Register("neoteqts4via6", setup)
}

// setup konfiguriert das Plugin in CoreDNS
func setup(c *caddy.Controller) error {
	p := NeoteqTS4via6{}
	c.OnStartup(func() error {
		fmt.Println("NeoteqTS4via6 Plugin geladen!")
		return nil
	})
	plugin.NextOrFailure(p.Name(), p.Next, context.Background(), nil, nil)
	return nil
}

// ServeDNS verarbeitet DNS-Anfragen für AAAA-Records
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

	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}

// Name gibt den Namen des Plugins zurück
func (p NeoteqTS4via6) Name() string {
	return "neoteqts4via6"
}
