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

	ret, err = queryByIP(host, service, proto)
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
