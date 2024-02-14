package validators

import (
	"testing"

	"github.com/maas/gomaasclient/entity"
	"github.com/stretchr/testify/assert"
)

type DummyMaaSAPIClient struct {
	images []entity.BootResource
}

func (d *DummyMaaSAPIClient) ListDNSServers() ([]entity.DNSResource, error) {
	return make([]entity.DNSResource, 0), nil
}

func (d *DummyMaaSAPIClient) ListOSImages() ([]entity.BootResource, error) {
	d.images = make([]entity.BootResource, 0)
	return d.images, nil
}

func TestFindingBootResources(t *testing.T) {
	maasRuleService := NewMaasRuleService(&DummyMaaSAPIClient{})

	images, _ := maasRuleService.ListOSImages()
	assert.Equal(t, len(images), 0)
}
