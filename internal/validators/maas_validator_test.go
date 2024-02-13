package validators

import (
	"github.com/maas/gomaasclient/entity"
)

type MockOSImageReader struct {
	images []entity.BootResource
}

type DummyMaaSAPIClient struct {
}

func (d *DummyMaaSAPIClient) ListDNSServers() ([]entity.DNSResource, error) {
	return make([]entity.DNSResource, 0), nil
}

func (d *DummyMaaSAPIClient) ListOSImages() ([]entity.BootResource, error) {
	return make([]entity.BootResource, 0), nil
}

func (m *MockOSImageReader) Get(params *entity.BootResourcesReadParams) ([]entity.BootResource, error) {
	return m.images, nil
}
