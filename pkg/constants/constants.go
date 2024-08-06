// Package constants defines the constants used in the plugin.
package constants

const (
	// PluginCode is the name of the plugin.
	PluginCode string = "MAAS"

	// ValidationTypeImage is the validation type for MAAS images
	ValidationTypeImage = "maas-image"

	// ValidationTypeUDNS is the validation type for MAAS upstream DNS
	ValidationTypeUDNS = "maas-upstream-dns"

	// ValidationTypeIDNS is the validation type for MAAS internal DNS
	ValidationTypeIDNS = "maas-internal-dns"

	// ValidationTypeResource is the validation type for MAAS resources
	ValidationTypeResource = "maas-resource"

	// ErrImageNotFound is the error message for when an image is not found
	ErrImageNotFound string = "failed to locate one or more image(s)"

	// ErrUDNSNotConfigured is the error message for when an upstream DNS server is not configured
	ErrUDNSNotConfigured string = "failed to locate one or more upstream DNS server(s)"

	// ErrIDNSNotConfigured is the error message for when an internal DNS server is not configured
	ErrIDNSNotConfigured string = "failed to locate one or more internal DNS server(s)"

	// ErrResourceNotFound is the error message for when a resource is not found
	ErrResourceNotFound string = "failed to locate one or more resource(s)"
)
