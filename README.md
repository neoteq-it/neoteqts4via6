# neoteqts4via6 CoreDNS Plugin

`neoteqts4via6` ist ein CoreDNS-Plugin zur Unterstützung der **Tailscale 4via6-Funktion**. Es ermöglicht die DNS-Auflösung von speziell formatierten Hostnamen in deterministisch generierte IPv6-Adressen, um überlappende IPv4-Netze im Tailscale-Mesh über IPv6 eindeutig erreichbar zu machen.

---

## Hintergrund: Tailscale 4via6

In komplexeren Tailscale-Setups, z. B. bei Standortvernetzungen oder hybriden Cloud-Umgebungen, kann es vorkommen, dass mehrere Netzwerke identische oder überlappende IPv4-Adressbereiche verwenden (z. B. `192.168.0.0/16`). Tailscale's **4via6**-Mechanismus löst dieses Problem, indem es jedem exportierten IPv4-Host eine eindeutige IPv6-Adresse zuweist, bestehend aus:

- der Original-IPv4-Adresse
- einer zusätzlichen eindeutigen **Verbindungs-ID**

Damit wird jede Quelle im Mesh eindeutig adressierbar – selbst wenn sich ihre IPv4-Adressen überschneiden.

Dieses Plugin stellt genau diese Funktionalität über DNS bereit.

---

## Funktionsweise

Das Plugin verarbeitet DNS-Anfragen mit folgendem Format: {ip}.via{id}.domain.tld.

### Beispiel:
192-168-1-100-via42.ts.net → wird aufgelöst zu: d7a:115c:a1e0:0b1a:0:2a:c0a8:0164

### Interpretation:

- **IPv4-Adresse**: `192.168.1.100`
- **Verbindungs-ID**: `42` → `0x2a`
- **Generierte IPv6-Adresse**: deterministisch, eindeutig und tailscale-kompatibel

---

## Aufbau der IPv6-Adresse

IPv6-Adressen werden nach folgendem Schema generiert:
fd7a:115c:a1e0:0b1a:0::

### Details:

| Feld       | Beschreibung                                      |
|------------|---------------------------------------------------|
| Präfix     | `fd7a:115c:a1e0:0b1a::/64` (Tailscale ULA-Space)  |
| ID         | Verbindungs-ID als 16-Bit-Hex                     |
| IPv4       | IPv4-Adresse in hexadezimaler Darstellung         |

Beispiel:

- ID `42` → `0x002a`
- IPv4 `192.168.1.100` → `c0a8:0164`
- Ergebnis: `fd7a:115c:a1e0:0b1a:0:2a:c0a8:0164`


