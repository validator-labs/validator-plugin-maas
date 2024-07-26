include build/makelib/common.mk
include build/makelib/plugin.mk

# Image URL to use all building/pushing image targets
IMG ?= quay.io/validator-labs/validator-plugin-maas:latest

# Helm vars
CHART_NAME=validator-plugin-maas

.PHONY: dev
dev:
	devspace dev -n validator

# Static Analysis / CI

chartCrds = chart/validator-plugin-maas/crds/validation.spectrocloud.labs_maasvalidators.yaml

reviewable-ext:
	rm $(chartCrds)
	cp config/crd/bases/validation.spectrocloud.labs_maasvalidators.yaml $(chartCrds)
