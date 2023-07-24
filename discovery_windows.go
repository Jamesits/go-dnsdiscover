package go_dnsdiscover

import (
	"github.com/yusufpapurcu/wmi"
)

type MSFT_DNSClientGlobalSetting struct {
	UseSuffixSearchList bool
	SuffixSearchList    []string
}

func DomainSearchList() ([]string, error) {
	var dst []MSFT_DNSClientGlobalSetting
	q := wmi.CreateQuery(&dst, "")
	err := wmi.QueryNamespace(q, &dst, "root\\StandardCimV2")
	if err != nil {
		return nil, err
	}
	if len(dst) == 0 {
		return nil, nil
	}
	if dst[0].UseSuffixSearchList == false {
		return []string{}, nil
	}
	return dst[0].SuffixSearchList, nil
}
