# Initial steps for code generation

```
export GOPATH=/Users/jaganaga/WA/golang
export ROOT_PACKAGE="github.com/dichque/grafana-operator"
export CUSTOM_RESOURCE_NAME="grafana"
export CUSTOM_RESOURCE_VERSION="v1"

mkdir -p pkg/apis/grafana/v1
```
- Create _pkg/apis/grafana/register.go_ , _pkg/apis/grafana/v1/types.go_ , _pkg/apis/grafana/v1/doc.go_ & _pkg/apis/grafana/v1/register.go_
- Commit & push it to the repository

```
go get -u k8s.io/code-generator/...
cd $GOPATH/src/k8s.io/code-generator

 ./generate-groups.sh all "$ROOT_PACKAGE/pkg/client" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION" --output-base "${GOPATH}/src" --go-header-file "hack/boilerplate.go.txt"
Generating deepcopy funcs
Generating clientset for grafana:v1 at github.com/dichque/grafana-operator/pkg/client/clientset
Generating listers for grafana:v1 at github.com/dichque/grafana-operator/pkg/client/listers
Generating informers for grafana:v1 at github.com/dichque/grafana-operator/pkg/client/informers

```

# Reference
- [Stringer Controller Development](https://medium.com/@trstringer/create-kubernetes-controllers-for-core-and-custom-resources-62fc35ad64a3)
- [Programming Kubernetes](https://github.com/programming-kubernetes/cnat/blob/master/cnat-client-go/pkg/apis/cnat/v1alpha1/types.go)
