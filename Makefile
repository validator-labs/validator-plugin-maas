include build/makelib/common.mk
include build/makelib/plugin.mk

# Image URL to use all building/pushing image targets
IMG ?= quay.io/validator-labs/validator-plugin-maas:latest

# Helm vars
CHART_NAME=validator-plugin-maas
