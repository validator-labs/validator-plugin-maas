package validators

import (
	"testing"

	"github.com/maas/gomaasclient/entity"
	"github.com/stretchr/testify/assert"
)

type DummyMaaSAPIClient struct {
}

func (d *DummyMaaSAPIClient) ListDNSServers() ([]entity.DNSResource, error) {
	return make([]entity.DNSResource, 0), nil
}

func (d *DummyMaaSAPIClient) ListOSImages() ([]entity.BootResource, error) {
	return make([]entity.BootResource, 0), nil
}

func TestFindingBootResources(t *testing.T) {
	maasRuleService := NewMaasRuleService(&DummyMaaSAPIClient{})

	images, _ := maasRuleService.ListOSImages()
	assert.Equal(t, len(images), 0)
}
