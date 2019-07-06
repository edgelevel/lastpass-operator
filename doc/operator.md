## operator-sdk

* [operator-sdk](https://github.com/operator-framework/operator-sdk)
* [A complete guide to Kubernetes Operator SDK](https://banzaicloud.com/blog/operator-sdk)
* [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions)
* [Groups and Versions and Kinds, oh my!](https://book.kubebuilder.io/cronjob-tutorial/gvks.html)
* [Controller Runtime Client API](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/client.md)
* [OperatorHub](https://operatorhub.io)

### Prerequisites

* [docker](https://docs.docker.com/install)
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

### Setup

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
