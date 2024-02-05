Hacking
=======

This project depends on the validator operator.

Before start hacking, install palette cli, and run:

```
palette validator install
```

Now, create the changes you need in the API. E.G add a new
custom rule definition in `api/v1alpha1/maasvalidator_types.go`

```
type MaasInstanceRule struct {
	// Host is a reference to the host URL of a MaaS instance
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

Build the binary with:
```
make
```

Install the CRDs with:
```
make install
```

Apply the samples:
```
kubectl apply -f config/samples/maasvalidator-osimages.yaml
```

Overriding target variables in the Makefile
-------------------------------------------

Many of the makefile targets have default variables defined. For example the default OCI image build
is defined as:

```
CONTAINER_TOOL ?= docker
```

You can override it via the command line with:

```
make docker-build CONTAINER_TOOL=podman
```

This will invoke `podman build` instead of `docker build`.
For a more permanent override you can use a `.env` file in the root of the project with:

```
$ cat .env
CONTAINER_TOOL=podman
```
This will override the default OCI build tool and use podman every time you invoke `make docker-build`.  
