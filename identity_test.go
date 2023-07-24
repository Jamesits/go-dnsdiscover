package go_dnsdiscover

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHostname(t *testing.T) {
	h, e := Hostname()
	assert.NoError(t, e)
	assert.NotEmpty(t, h)
	t.Logf("hostname = %v", h)
}

func TestHostnameFQDN(t *testing.T) {
	h, e := HostnameFQDN()
	assert.NoError(t, e)
	assert.NotEmpty(t, h)
	assert.True(t, strings.Contains(h, "."))
	t.Logf("FQDN = %v", h)
}

func TestHostnameShort(t *testing.T) {
	h, e := HostnameShort()
	assert.NoError(t, e)
	assert.NotEmpty(t, h)
	assert.False(t, strings.Contains(h, "."))
	t.Logf("short hostname = %v", h)
}

func TestHostnameLinkLocal(t *testing.T) {
	h, e := HostnameLinkLocal()
	assert.NoError(t, e)
	assert.NotEmpty(t, h)
	assert.False(t, strings.Contains(h, "."))
	t.Logf("link local hostname = %v", h)
}
