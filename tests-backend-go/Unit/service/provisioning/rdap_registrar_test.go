package provisioning_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRDAPRegistrarDriver(t *testing.T) {
	driver := provisioning.NewRDAPRegistrarDriver()
	ctx := context.Background()

	t.Run("Check availability of known taken domain (cache & live lookup)", func(t *testing.T) {
		res, err := driver.CheckAvailability(ctx, "google.com")
		require.NoError(t, err)
		assert.Equal(t, "google.com", res.DomainName)
		assert.Equal(t, "com", res.TLD)
		assert.False(t, res.IsAvailable)
		assert.Equal(t, int64(129900), res.Price)
	})

	t.Run("Check availability of unique unregistered domain format", func(t *testing.T) {
		res, err := driver.CheckAvailability(ctx, "foss-test-unregistered-domain-99419.com")
		require.NoError(t, err)
		assert.Equal(t, "foss-test-unregistered-domain-99419.com", res.DomainName)
		assert.Equal(t, "com", res.TLD)
	})

	t.Run("Register and check taken state transition", func(t *testing.T) {
		domainName := "brandnewdomain2026.id"
		reg, err := driver.RegisterDomain(ctx, provisioning.DomainRegistrationRequest{
			DomainName: domainName,
			Years:      1,
			Nameservers: []string{"ns1.customdns.org", "ns2.customdns.org"},
		})
		require.NoError(t, err)
		assert.Equal(t, "active", reg.Status)
		assert.Contains(t, reg.AuthCode, "EPP-")

		// Re-checking availability should now return taken
		res, err := driver.CheckAvailability(ctx, domainName)
		require.NoError(t, err)
		assert.False(t, res.IsAvailable)
	})

	t.Run("Invalid domain format returns error", func(t *testing.T) {
		_, err := driver.CheckAvailability(ctx, "invalid_domain")
		require.Error(t, err)
	})
}
