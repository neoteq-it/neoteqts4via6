package neoteqts4via6

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/coredns/caddy"
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
	p := NeoteqTS4via6{}
	c.OnStartup(func() error {
		fmt.Println("NeoteqTS4via6 Plugin geladen!")
		return nil
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

	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}

func (p NeoteqTS4via6) Name() string {
	return "neoteqts4via6"
}

func ResolveIPv6(query string) (string, error) {
	parts := strings.Split(query, ".")
	if len(parts) < 3 {
		return "", fmt.Errorf("ungültige Anfrage")
	}

	ipv4Str := strings.ReplaceAll(parts[0], "-", ".")
	idStr := strings.TrimPrefix(parts[1], "via")

	ipv4Parts := strings.Split(ipv4Str, ".")
	if len(ipv4Parts) != 4 {
		return "", fmt.Errorf("ungültige IPv4-Adresse")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return "", fmt.Errorf("ungültige ID")
	}

	ipv4Bytes := make([]int, 4)
	for i := 0; i < 4; i++ {
		ipv4Bytes[i], err = strconv.Atoi(ipv4Parts[i])
		if err != nil || ipv4Bytes[i] < 0 || ipv4Bytes[i] > 255 {
			return "", fmt.Errorf("ungültige IPv4-Adresse")
		}
	}

	ipv6 := fmt.Sprintf("fd7a:115c:a1e0:b1a:0:%x:%02x%02x:%02x%02x",
		id, ipv4Bytes[0], ipv4Bytes[1], ipv4Bytes[2], ipv4Bytes[3])

	return ipv6, nil
}
