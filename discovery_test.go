package go_dnsdiscover

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestDomain string

func TestMain(m *testing.M) {
	// assuming we have an AD deployment here...
	l, _ := DomainSearchList()
	TestDomain = l[0]
	m.Run()
}

func TestGetDomainSearchList(t *testing.T) {
	l, err := DomainSearchList()
	assert.NoError(t, err)
	assert.NotZero(t, len(l))
	t.Logf("Search list: %v\n", l)
}

func TestQueryByIP(t *testing.T) {
	srv, err := queryByIP(TestDomain, "ldap", "tcp")
	assert.NoError(t, err)
	assert.NotZero(t, len(srv))
	t.Logf("Services: \n")
	for _, s := range srv {
		t.Logf("\t%v\n", s)
	}
}

func TestQueryBySRV(t *testing.T) {
	srv, err := queryBySRV(TestDomain, "ldap", "tcp")
	assert.NoError(t, err)
	assert.NotZero(t, len(srv))
	t.Logf("Services: \n")
	for _, s := range srv {
		t.Logf("\t%v\n", s)
	}
}

func TestQuerySingleService(t *testing.T) {
	srv, err := querySingleServiceEndpoint(TestDomain, "kerberos", "tcp")
	assert.NoError(t, err)
	assert.NotZero(t, len(srv))
	t.Logf("Services: \n")
	for _, s := range srv {
		t.Logf("\t%v\n", s)
	}
}

func TestQueryServiceEndpoint(t *testing.T) {
	srv, err := QueryServiceEndpoint(TestDomain, "kerberos", "tcp")
	assert.NoError(t, err)
	assert.NotZero(t, len(srv))
	t.Logf("Services: \n")
	for _, s := range srv {
		t.Logf("\t%v\n", s)
	}
}

func TestQueryServiceEndpointDefault(t *testing.T) {
	srv, err := QueryServiceEndpoint("", "kerberos", "tcp")
	assert.NoError(t, err)
	assert.NotZero(t, len(srv))
	t.Logf("Services: \n")
	for _, s := range srv {
		t.Logf("\t%v\n", s)
	}
}

func TestQueryServiceEndpointSanity(t *testing.T) {
	// 2b2t utilizes SRV record
	srv, err := QueryServiceEndpoint("2b2t.org", "minecraft", "tcp")
	assert.NoError(t, err)
	assert.NotZero(t, len(srv))
	assert.EqualValues(t, "connect.2b2t.org.", srv[0].Target)
	t.Logf("Services: \n")
	for _, s := range srv {
		t.Logf("\t%v\n", s)
	}

	// google.com uses an A or AAAA record
	srv, err = QueryServiceEndpoint("google.com", "https", "tcp")
	assert.NoError(t, err)
	assert.NotZero(t, len(srv))
	t.Logf("Services: \n")
	for _, s := range srv {
		t.Logf("\t%v\n", s)
	}
}
