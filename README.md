# lastpass-operator

> TODO

* [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions)
* [Groups and Versions and Kinds, oh my!](https://book.kubebuilder.io/cronjob-tutorial/gvks.html)
* [Controller Runtime Client API](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/client.md)
* [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret)
* [OperatorHub](https://operatorhub.io)

## Development

### LastPass

* [lastpass-cli](https://github.com/lastpass/lastpass-cli)

CLI examples
```bash
# login
echo PWD | LPASS_DISABLE_PINENTRY=1 lpass login --trust EMAIL

# list
lpass ls

# retrieve passwords
lpass show <GROUP>/<NAME> --json --expand-multi

# logout
lpass logout --force
```

Docker
```bash
# build image
docker build -t niqdev/lastpass-cli .

# temporary container
docker run --rm -it niqdev/lastpass-cli /bin/sh
lpass --version

# execute command inline
docker run --rm -it niqdev/lastpass-cli lpass --version
```

### operator-sdk

* [A complete guide to Kubernetes Operator SDK](https://banzaicloud.com/blog/operator-sdk)
* [operator-sdk](https://github.com/operator-framework/operator-sdk)

**Prerequisites**

* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)
* [helm](https://helm.sh/docs/using_helm/#installing-helm)
* [go](https://golang.org/doc)
* [dep](https://golang.github.io/dep/docs/introduction.html)
* [operator-sdk](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md)

Ubuntu
```bash
# kubectl
sudo snap install kubectl --classic

# helm
sudo snap install helm --classic

# go
sudo snap install --classic go

# dep
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# operator-sdk (from source)
go get -d github.com/operator-framework/operator-sdk
cd $GOPATH/src/github.com/operator-framework/operator-sdk
git checkout master
make dep
make install
```

macOS
```bash
# kubectl
brew install kubernetes-cli

# helm
brew install kubernetes-helm

# go
brew install go

# dep
brew install dep

# operator-sdk
brew install operator-sdk
```

Setup go [workspace](https://golang.org/doc/code.html#Workspaces)
```bash
# add to .bashrc or .bash_profile
export GOPATH=$HOME/go
export PATH=$PATH:$(go env GOPATH)/bin
```

Initialize project
```bash
# create project
mkdir -p $GOPATH/src/github.com/niqdev && cd $_
operator-sdk new lastpass-operator --dep-manager=dep

# add crd
operator-sdk add api --api-version=niqdev.com/v1alpha1 --kind=LastPassSecret
operator-sdk generate k8s

# add controller
operator-sdk add controller --api-version=niqdev.com/v1alpha1 --kind=LastPassSecret
```

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

1) apply master secret manually
2) operator: check wait for the master secret to be available
3) retrieve master credentials - no MFA allowed or request login allowed (skip secret)
4) exec `lpass` login
5) create secrets
6) add status: verify `last_modified_gmt` and `last_touch`
7) logout every time after all?

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
