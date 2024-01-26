Hacking
=======

This project depends on the validator operator.

Before you start hacking, install palette cli, and run:

```
palette validator install
```

Now, create the changes you need in the API. E.G add a new
custom rule definition in `api/v1alpha1/maasvalidator_types.go`

```
type MaasInstanceRule struct {
	// Host is a reference to the host URL of an OCI compliant registry
	Host string `json:"host" yaml:"host"`

	// OSImages is a list of bootable os images
	OSImages []OSImage `json:"bootable-images,omitempty" yaml:"artifacts,omitempty"`

	// Auth provides authentication information for the MaaS Instance
	Auth Auth `json:"auth,omitempty" yaml:"auth,omitempty"`
}
```

Now you need to re-create the client code:

```
make generate
```
This will update: `api/v1alpha1/zz_generated.deepcopy.go`

```
make manifests
```

This will update: `config/crd/bases/validation.spectrocloud.labs_maasvalidators.yaml`
