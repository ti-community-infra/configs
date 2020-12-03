module github.com/tidb-community-bots/configs

go 1.15

require (
	github.com/ghodss/yaml v1.0.0
	github.com/sirupsen/logrus v1.7.0
	k8s.io/apimachinery v0.19.3
	k8s.io/test-infra v0.0.0-20201114015505-f09ff0e80535
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.19.3
	k8s.io/client-go => k8s.io/client-go v0.19.3
)
