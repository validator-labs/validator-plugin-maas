// Package maas provides a client for the MaaS API
package maas

import (
	gomaasclient "github.com/canonical/gomaasclient/client"

	"github.com/canonical/gomaasclient/entity"
)

// API is a struct which contains the Maas client
type API struct {
	ImageClient
	DNSServerClient
	ResourceClient
}

// ImageClient is a struct which contains the Maas client for OS images
type ImageClient struct {
	Client *gomaasclient.Client
}

// DNSServerClient is a struct which contains the Maas client for DNS servers
type DNSServerClient struct {
	Client *gomaasclient.Client
}

// ResourceClient is a struct which contains the Maas client for compute resources
type ResourceClient struct {
	Client *gomaasclient.Client
}

// NewAPI returns a new API for MaaS clients
func NewAPI(client *gomaasclient.Client) *API {
	return &API{
		ImageClient:     ImageClient{Client: client},
		DNSServerClient: DNSServerClient{Client: client},
		ResourceClient:  ResourceClient{Client: client},
	}
}

// ListOSImages returns a list of OS images from the Maas API
func (i *ImageClient) ListOSImages() ([]entity.BootResource, error) {
	if i.Client != nil {
		images, err := i.Client.BootResources.Get(&entity.BootResourcesReadParams{})
		if err != nil {
			return make([]entity.BootResource, 0), err
		}
		return images, nil
	}
	return make([]entity.BootResource, 0), nil
}

// ListDNSServers returns a list of DNS servers from the Maas API
func (d *DNSServerClient) ListDNSServers() ([]entity.DNSResource, error) {
	if d.Client != nil {
		dnsresources, err := d.Client.DNSResources.Get()
		if err != nil {
			return make([]entity.DNSResource, 0), err
		}
		return dnsresources, nil
	}
	return make([]entity.DNSResource, 0), nil
}
