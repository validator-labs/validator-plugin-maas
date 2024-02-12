package validators

import "github.com/maas/gomaasclient/entity"

type MockOSImageReader struct {
	images []entity.BootResource
}

func (m *MockOSImageReader) Get(params *entity.BootResourcesReadParams) ([]entity.BootResource, error) {
	return m.images, nil
}
