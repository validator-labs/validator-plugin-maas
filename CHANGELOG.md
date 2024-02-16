# Changelog

## [0.0.1](https://github.com/spectrocloud-labs/validator-plugin-maas/compare/v0.0.2...v0.0.1) (2024-02-16)


### Features

* add Helm chart ([#25](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/25)) ([f4295ae](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/f4295ae9a509c52763c12ba01458d8d0150b0bae))
* allow initContainer image to be passed in via values.yaml ([#27](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/27)) ([50c8647](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/50c8647f76cc70453b1ec1a5f7e307fcda839235))
* implement OCI registry validation spec ([#6](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/6)) ([f62c494](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/f62c494d3a44bcf99c9d0bccecd1af2b8bc3ae78))
* support validating list of oci artifacts ([#16](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/16)) ([d0cbecc](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/d0cbecc24614a9a6ddf2a34e71e01ce23a313d8c))


### Bug Fixes

* CRD validation for rule host uniqueness ([#56](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/56)) ([8dbdc15](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/8dbdc15a2d23225e94630771d26eab26439c721c))
* **deps:** update aws-sdk-go-v2 monorepo ([#55](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/55)) ([af7f8a4](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/af7f8a47423f262b9d491b9bac7bda3ba8c21ac8))
* **deps:** update aws-sdk-go-v2 monorepo ([#61](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/61)) ([b733807](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/b7338076acaac8d86965a826cc2b27b1c626390b))
* **deps:** update aws-sdk-go-v2 monorepo ([#67](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/67)) ([c1c5d0e](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/c1c5d0e2543c53e4d36c8affb484c0847c6d6275))
* **deps:** update aws-sdk-go-v2 monorepo ([#76](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/76)) ([55d84a8](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/55d84a85bd44a435371068ea0329455332008bee))
* **deps:** update aws-sdk-go-v2 monorepo ([#81](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/81)) ([1b4d64d](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/1b4d64d305faaa6008e29d7d7723f3acf9188029))
* **deps:** update kubernetes packages to v0.28.4 ([#17](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/17)) ([f346f63](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/f346f631c50d2fdc6236603055c792f116c554df))
* **deps:** update kubernetes packages to v0.29.2 ([#6](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/6)) ([2e3014a](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/2e3014a3720ed3476b308c82a89e051226df51f2))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.10 ([#44](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/44)) ([4e221d7](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/4e221d7cd8868c4677a84daaee5e1382c20f3ba5))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.11 ([#48](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/48)) ([8567bef](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/8567bef0cf026b79fb1fff133eb90da9d351d652))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.6 ([#28](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/28)) ([b12dabe](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/b12dabe9730e9e12a48e979f796fde71dbd551a0))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.8 ([#32](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/32)) ([3eb0824](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/3eb08241bd645d73cd50182fd562f941171b4a30))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.25.9 ([#37](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/37)) ([74b6eae](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/74b6eae6b16df62c9866898b76d0ceaf95edbe4f))
* **deps:** update module github.com/aws/aws-sdk-go-v2/config to v1.26.6 ([#85](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/85)) ([939b7cc](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/939b7ccd1d5ec2d9be9649b9fb0c51f449e57b24))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.23.2 ([#30](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/30)) ([6375b2b](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/6375b2bafbdcaa8649691eaba15ba52ec8eb80d9))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.23.3 ([#34](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/34)) ([a55ada3](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/a55ada393e0ef05510388a152b4e0a03a573d3d4))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.24.0 ([#39](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/39)) ([d6d4314](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/d6d4314c9e21a89c22542f099117bbb7ab20a5e4))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.24.1 ([#45](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/45)) ([9ba631c](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/9ba631cc4b624cbfd8941d64b451a61b77392775))
* **deps:** update module github.com/aws/aws-sdk-go-v2/service/ecr to v1.24.2 ([#47](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/47)) ([7275869](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/727586939c110c34c2c60786279356ae40c3131b))
* **deps:** update module github.com/google/go-containerregistry to v0.17.0 ([#42](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/42)) ([a6c2c58](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/a6c2c582ffa74a0c97a87f318f2b763029725037))
* **deps:** update module github.com/google/go-containerregistry to v0.18.0 ([#77](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/77)) ([bfe4961](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/bfe49617b3649e86f50a1062a5b8d923dade8f7c))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.1 ([#15](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/15)) ([23673ac](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/23673ac0092fac7eecc78f7d92c249b385537c39))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.2 ([#35](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/35)) ([4f44a26](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/4f44a26a67d141a8b953cd937d07c1c0482087eb))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.14.0 ([#73](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/73)) ([f340000](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/f34000094b81b712a86ced4b28b21aa2c7d295d9))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.15.0 ([#78](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/78)) ([d0599fb](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/d0599fbd0facae1048feedfabb037d564c37fc72))
* **deps:** update module github.com/onsi/gomega to v1.31.0 ([#79](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/79)) ([d6d17a9](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/d6d17a92de92d9167533b354f3f5bacac5ab033c))
* **deps:** update module github.com/onsi/gomega to v1.31.1 ([#83](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/83)) ([731a624](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/731a6241b95df5dea83a8bac65b93801ebf3c25c))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.18 ([#14](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/14)) ([58c78f4](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/58c78f43a7f21d8d22381042ea73fbe5f3b7f0d0))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.21 ([#18](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/18)) ([9373c1d](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/9373c1d3541397948eca4a93df61c3a628661b56))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.25 ([#21](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/21)) ([76f1b24](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/76f1b247a7bf69d990a0539e8ec73260cfe7ad5a))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.26 ([#36](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/36)) ([2c18421](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/2c184212663c4ac97048b50fb2399f163f519036))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.27 ([#43](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/43)) ([1113a9e](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/1113a9ea320a34de327970093542681998959669))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.28 ([#52](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/52)) ([4fb5e57](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/4fb5e57a0065602ef699b03c6c72928d246dbfb0))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.30 ([#66](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/66)) ([dfc8fd7](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/dfc8fd7eeddd644eec71ebc9402940f6bfbc8b98))
* **deps:** update module github.com/spectrocloud-labs/validator to v0.0.32 ([#69](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/69)) ([105f4ce](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/105f4ce027834c243998926dbd2c2df7646daf43))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.2 ([#7](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/7)) ([f609e6b](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/f609e6b01d08fd496f40b54094f0415992e04d06))
* ensure codecov is run when code is pushed to main ([#59](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/59)) ([22da463](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/22da46306fc2d96f6aa10baaad0ff347a4ceb139))
* fix link to oci issues in readme ([#41](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/41)) ([b3c1cea](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/b3c1cea77e46bab05727be8bca15b85f2687c6e4))
* set owner references on validation result to ensure cleanup ([#19](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/19)) ([9c7c28d](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/9c7c28d1e69b9488263537e48415818826d96ebf))
* update leader election id ([#46](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/46)) ([976487b](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/976487bfe8edaebe638d3cf067f787d5ec2385b0))


### Other

* add license badge ([1eb5f1b](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/1eb5f1b2ceafc7656816f42b4f51c11ad0057aba))
* **deps:** pin dependencies ([#9](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/9)) ([9876cd7](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/9876cd701be178016231d02661a78db1f2f48c85))
* **deps:** update actions/checkout action to v4 ([#10](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/10)) ([cd110af](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/cd110af99d4eed651d89dabe5565bcedcb3f4c35))
* **deps:** update actions/setup-python action to v5 ([#53](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/53)) ([2afc37b](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/2afc37bda4e3f6e0e23ab57938f33af4ada2cd59))
* **deps:** update actions/upload-artifact action to v4 ([#63](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/63)) ([4fbdafc](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/4fbdafc6dbe48a19049667f02860815b469c83fe))
* **deps:** update actions/upload-artifact digest to 1eb3cb2 ([#74](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/74)) ([84e1f0a](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/84e1f0a1cc3a1b51785277c2333125c14ead1c2f))
* **deps:** update actions/upload-artifact digest to 694cdab ([#82](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/82)) ([7518f3f](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/7518f3f94357116c9998423da3c5efaae4d8da5e))
* **deps:** update anchore/sbom-action action to v0.15.0 ([#23](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/23)) ([34253f0](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/34253f03e491ebecc0ce8631d56558cc16bb4b82))
* **deps:** update anchore/sbom-action action to v0.15.1 ([#51](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/51)) ([ebf1d17](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/ebf1d1770babd60e53a642041fe5d026f56c8838))
* **deps:** update anchore/sbom-action action to v0.15.2 ([#70](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/70)) ([8e4c2f8](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/8e4c2f82ffde501af505b08404a65927754b859d))
* **deps:** update anchore/sbom-action action to v0.15.3 ([#71](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/71)) ([0e3dea8](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/0e3dea878db1a962d81e11effd0737fb30d95eaf))
* **deps:** update anchore/sbom-action action to v0.15.4 ([#80](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/80)) ([e1771fa](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/e1771fa7d145004931095b3c68ff9ce424b388cc))
* **deps:** update anchore/sbom-action action to v0.15.5 ([#84](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/84)) ([9576d15](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/9576d15e3ac41592bb54d660e0fbdad655bb1439))
* **deps:** update anchore/sbom-action action to v0.15.8 ([#5](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/5)) ([1e19586](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/1e19586dfb7a04df2ceb1c0fb2237ba0b2cfe1b7))
* **deps:** update docker/build-push-action digest to 4a13e50 ([#20](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/20)) ([eace63e](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/eace63e7d49fc14c8d1f8d0427bd11039bef140d))
* **deps:** update gcr.io/spectro-images-public/golang docker tag to v1.22 ([#16](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/16)) ([8f4819a](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/8f4819ad4d3e4dd64dfcdaf4b4c2d46c56b784f8))
* **deps:** update gcr.io/spectro-images-public/golang docker tag to v1.22 ([#72](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/72)) ([ff3312c](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/ff3312cd60da5cdb3aa0ee73ad1903ba5dff9a3c))
* **deps:** update golang docker tag to v1.22 ([#11](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/11)) ([824ca8b](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/824ca8b9fa48007f6835097f5d59948d873a889e))
* **deps:** update google-github-actions/release-please-action action to v4 ([#50](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/50)) ([78956e2](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/78956e209a28d7cb4d756e360ecad096c3c48576))
* **deps:** update google-github-actions/release-please-action digest to a2d8d68 ([#58](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/58)) ([23821e5](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/23821e5a4388b5e519383c385ab24e75f9210893))
* **deps:** update google-github-actions/release-please-action digest to cc61a07 ([#64](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/64)) ([b1dee9b](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/b1dee9bf69fd3ed5114a8c330711c69eec0ce913))
* fix platform specification for manager image ([#13](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/13)) ([539e8be](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/539e8be372a623125d1ed04e602833c59acddd93))
* **main:** release 0.0.1 ([#38](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/38)) ([cdd256d](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/cdd256db984995f939981fc70514a1953d829415))
* **main:** release 0.0.2 ([#40](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/40)) ([53db66e](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/53db66e0c57a76b063011e7349c326fdd855fb74))
* release 0.0.1 ([d3ecdef](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/d3ecdef9321650d8e947e058a2967d622bf0dad6))
* release 0.0.1 ([d4b32af](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/d4b32afa1737b2ca4dc39e907bbb4bee871e15fc))
* specify platform in Dockerfile and docker-build make target ([#12](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/12)) ([a88d182](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/a88d1820503bdfc6c2f99690db2a0bcd6befc5dc))
* switch back to public bulwark images ([010e7f8](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/010e7f842a54cb0f9e0f572618007ad85009f766))
* update spectrocloud-labs/validator dependency to v0.0.15 ([f62c494](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/f62c494d3a44bcf99c9d0bccecd1af2b8bc3ae78))


### Refactoring

* switch from oras to go-containerregistry ([#24](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/24)) ([eef0013](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/eef0013a7d1072f55bb3356304f287ad1cc61ff4))
* switch init container to image with ca-certificates pre installed ([#33](https://github.com/spectrocloud-labs/validator-plugin-maas/issues/33)) ([4550f4b](https://github.com/spectrocloud-labs/validator-plugin-maas/commit/4550f4bedb9807d8578fcc56d7fc4e3309cd6d8b))

## [0.0.0] Initial work based on OCI plugin
