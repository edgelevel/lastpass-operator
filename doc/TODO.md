## Development

Run on [minkube](https://github.com/kubernetes/minikube)
```bash
# requires virtualbox
minikube start

# apply crd
kubectl apply -f deploy/crds/niqdev_v1alpha1_lastpasssecret_crd.yaml

# run locally
export OPERATOR_NAME=lastpass-operator
operator-sdk up local --namespace=default --verbose
```

---

```bash
kubectl apply -f lastpass-master-secret.yaml
kubectl get secrets
kubectl get secret lastpass-master-secret -o yaml
echo '' | base64 --decode
```

---

**TODO**

```bash
# clone sources
mkdir -p $GOPATH/src/github.com/niqdev && cd $_
git clone git@github.com:niqdev/lastpass-operator.git

# first time only
dep init

# add dependencies
dep ensure -add github.com/USER/DEP1 github.com/USER/DEP2
# example
dep ensure -add github.com/spf13/cobra

# verify and update all dependencies
dep status
dep check
dep ensure -update

# resolve dependencies
dep ensure

# init cli
cobra init . --pkg-name lastpass-operator

# run
go run main.go

# compile
go build $GOPATH/src/github.com/niqdev/lastpass-operator

# compile and build executable
go install $GOPATH/src/github.com/niqdev/lastpass-operator

# test
go test $GOPATH/src/github.com/niqdev/lastpass-operator
```
