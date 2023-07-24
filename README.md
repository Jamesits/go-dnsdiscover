# go-dnsdiscover

[![Go Reference](https://pkg.go.dev/badge/github.com/Jamesits/go-dnsdiscover.svg)](https://pkg.go.dev/github.com/Jamesits/go-dnsdiscover)

DNS-based service discovery support library.

Provides:
- Hostname retrieval
  - short form
  - FQDN/full form
- SRV-based service lookup, with A/AAAA based fallback

TODO:
- Provide a custom `*net.Resolver` - https://github.com/golang/go/issues/12503
