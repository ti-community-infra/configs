module github.com/tidb-community-bots/prow-configs

go 1.15

require (
	github.com/sirupsen/logrus v1.7.0
	k8s.io/test-infra v0.0.0-20201114015505-f09ff0e80535
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/api => k8s.io/api v0.19.3
	k8s.io/client-go => k8s.io/client-go v0.19.3
)
