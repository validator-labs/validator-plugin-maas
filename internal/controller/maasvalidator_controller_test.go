package controller

import (
	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"
	//+kubebuilder:scaffold:imports
)

const MaasValidatorName = "maas-validator"

var _ = Describe("MaaSValidator controller", Ordered, func() {

	BeforeEach(func() {
		// toggle true/false to enable/disable the OCIValidator controller specs
		if false {
			Skip("skipping")
		}
	})
})
