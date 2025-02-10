package neoteqts4via6

import (
	"fmt"
	"strconv"
	"strings"
)

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
