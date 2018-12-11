# Sec § — a golang opiniated dependency updater

```shell
$ sec
 INFO Gather dependencies to update…
 WARN Skip dependency github.com/docker/distribution : constraint branch openshift-3.10-docker-edc3ab2 (override)
 WARN Skip dependency github.com/golang/glog : constraint branch master
 WARN Skip dependency github.com/google/gofuzz : constraint branch master
 WARN Skip dependency github.com/openshift/api : constraint branch release-3.11 (override)
 WARN Skip dependency github.com/openshift/client-go : constraint branch release-3.11
 WARN Skip dependency github.com/openshift/library-go : constraint branch release-3.11 (override)
 WARN Skip dependency github.com/openshift/origin : constraint branch release-3.11 (override)
 WARN Skip dependency golang.org/x/crypto : constraint branch master
 WARN Skip dependency golang.org/x/net : constraint branch master
 WARN Skip dependency golang.org/x/sys : constraint branch master
 WARN Skip dependency golang.org/x/time : constraint branch master
 WARN Skip dependency k8s.io/api : constraint branch origin-3.11-kubernetes-1.11.1 (override)
 WARN Skip dependency k8s.io/apimachinery : constraint branch origin-3.11-kubernetes-1.11.1 (override)
 WARN Skip dependency k8s.io/client-go : constraint branch origin-3.11-kubernetes-1.11.1 (override)
 WARN Skip dependency k8s.io/kubernetes : constraint branch origin-3.11-kubernetes-1.11.1 (override)
 INFO Update dependency github.com/ghodss/yaml…
 INFO … no updates for github.com/ghodss/yaml
 INFO Update dependency github.com/gogo/protobuf…
 INFO … - github.com/gogo/protobuf : v1.1.1 -> v1.2.0
 INFO Execute command `go test ./...` (timeout 5m0s)…
 INFO Commit dependency changes for github.com/gogo/protobuf…
 INFO Update dependency github.com/imdario/mergo…
 INFO … no updates for github.com/imdario/mergo
 INFO Update dependency github.com/inconshreveable/mousetrap…
 INFO … no updates for github.com/inconshreveable/mousetrap
 INFO Update dependency github.com/json-iterator/go…
 INFO … no updates for github.com/json-iterator/go
 INFO Update dependency github.com/modern-go/concurrent…
 INFO … no updates for github.com/modern-go/concurrent
 INFO Update dependency github.com/modern-go/reflect2…
 INFO … no updates for github.com/modern-go/reflect2
 INFO Update dependency github.com/opencontainers/go-digest…
 INFO … no updates for github.com/opencontainers/go-digest
 INFO Update dependency github.com/openshift/source-to-image…
 INFO … - github.com/openshift/source-to-image : v1.1.11 -> v1.1.12
 INFO Execute command `go test ./...` (timeout 5m0s)…
 INFO Commit dependency changes for github.com/openshift/source-to-image…
 INFO Update dependency github.com/pkg/errors…
 INFO … no updates for github.com/pkg/errors
 INFO Update dependency github.com/spf13/cobra…
 INFO … no updates for github.com/spf13/cobra
 INFO Update dependency github.com/spf13/pflag…
 INFO … no updates for github.com/spf13/pflag
 INFO Update dependency golang.org/x/text…
 INFO … no updates for golang.org/x/text
 INFO Update dependency gopkg.in/inf.v0…
 INFO … no updates for gopkg.in/inf.v0
 INFO Update dependency gopkg.in/yaml.v2…
 INFO … - gopkg.in/yaml.v2 : v2.2.1 -> v2.2.2
 INFO Execute command `go test ./...` (timeout 5m0s)…
 INFO Commit dependency changes for gopkg.in/yaml.v2…
$ git log
* 97f9225 - (HEAD -> master) gopkg.in/yaml.v2: v2.2.1 -> v2.2.2 (15 seconds ago)
* feed3e9 - github.com/openshift/source-to-image: v1.1.11 -> v1.1.12 (70 seconds ago)
* 79d4a83 - github.com/gogo/protobuf: v1.1.1 -> v1.2.0 (2 minutes ago)
* cad8640 - (origin/master) Add apache v2 license 👅 (3 weeks ago)
```
