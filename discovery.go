package go_dnsdiscover

import (
	"net"
	"strings"
)

var portList = map[string]map[string]int{
	"tcp": {
		"minecraft": 25565,
	},
}

// AddPortAssignment assigns a service with a port.
// Normally you should edit `/etc/services` for this, but there are cases where editing that file is impossible, then
// you can do it here.
// proto: "tcp" or "udp"
// service: the name of the service, e.g. "https" or "ldap" or "minecraft".
func AddPortAssignment(proto string, service string, port int) error {
	if _, ok := portList[proto]; !ok {
		portList[proto] = make(map[string]int)
	}

	portList[proto][service] = port
	return nil
}

// LookupPort returns the port of a service.
// proto: "tcp" or "udp"
// service: the name of the service, e.g. "https" or "ldap"
func LookupPort(proto string, service string) (port int, err error) {
	port, err = net.LookupPort(proto, service)
	if err == nil {
		return
	}

	if p, ok := portList[proto]; ok {
		if port, ok := p[service]; ok {
			return port, nil
		}
	}

	return 0, err
}

func queryByIP(host string, service string, proto string) (ret []*net.SRV, err error) {
	ips, err := net.LookupIP(host)
	if ips == nil || len(ips) == 0 {
		return nil, err
	}

	port, err := LookupPort(proto, service)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		ret = append(ret, &net.SRV{
			Target:   ip.String(),
			Port:     uint16(port),
			Priority: 0,
			Weight:   0,
		})
	}
	return ret, nil
}

func fakeSRV(host string, service string, proto string) (ret []*net.SRV, err error) {
	// try resolve it first
	ips, err := net.LookupIP(host)
	if ips == nil || len(ips) == 0 {
		return nil, err
	}

	port, err := LookupPort(proto, service)
	if err != nil {
		return nil, err
	}

	ret = append(ret, &net.SRV{
		Target:   host,
		Port:     uint16(port),
		Priority: 0,
		Weight:   0,
	})
	return
}

func queryBySRV(host string, service string, proto string) (ret []*net.SRV, err error) {
	// https://www.ietf.org/rfc/rfc2782.txt
	_, addrs, err := net.LookupSRV(service, proto, host)
	return addrs, err
}

func querySingleServiceEndpoint(host string, service string, proto string) (ret []*net.SRV, err error) {
	ret, err = queryBySRV(host, service, proto)
	if err == nil {
		return
	}

	ret, err = fakeSRV(host, service, proto)
	if err == nil {
		return
	}

	return nil, err
}

// QueryServiceEndpoint returns a list of possible endpoints of a service available under the host.
func QueryServiceEndpoint(host string, service string, proto string) (ret []*net.SRV, err error) {
	if !strings.Contains(host, ".") {
		// process Domain Search Option https://www.rfc-editor.org/rfc/rfc3397
		var ds []string
		ds, err = DomainSearchList()
		if err != nil {
			goto direct
		}

		for _, d := range ds {
			var fqdn string
			if host != "" {
				fqdn = strings.Join([]string{host, d}, ".")
			} else {
				fqdn = d
			}

			ret, err = querySingleServiceEndpoint(fqdn, service, proto)
			if err == nil {
				return ret, nil
			}
		}
	}

direct:
	return querySingleServiceEndpoint(host, service, proto)
}
