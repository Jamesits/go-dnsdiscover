package go_dnsdiscover

import (
	"golang.org/x/sys/windows"
	"os"
	"syscall"
)

func hostname(format int) (string, error) {
	n := uint32(64)
	for {
		b := make([]uint16, n)
		err := windows.GetComputerNameEx(uint32(format), &b[0], &n)
		if err == nil {
			return syscall.UTF16ToString(b[:n]), nil
		}
		if err != syscall.ERROR_MORE_DATA {
			return "", os.NewSyscallError("ComputerNameEx", err)
		}

		// If we received an ERROR_MORE_DATA, but n doesn't get larger,
		// something has gone wrong, and we may be in an infinite loop
		if n <= uint32(len(b)) {
			return "", os.NewSyscallError("ComputerNameEx", err)
		}
	}
}

// HostnameFQDN returns the FQDN of the current host (hostname.domain).
func HostnameFQDN() (string, error) {
	return hostname(windows.ComputerNamePhysicalDnsFullyQualified)
}

// HostnameShort returns the hostname part of the hostname,
func HostnameShort() (string, error) {
	return hostname(windows.ComputerNamePhysicalDnsHostname)
}

// HostnameLinkLocal returns the hostname used for link-local service discovery.
func HostnameLinkLocal() (string, error) {
	return hostname(windows.ComputerNameNetBIOS)
}
