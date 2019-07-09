## Setup

* [docker](https://docs.docker.com/install)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)
* [helm](https://helm.sh/docs/using_helm/#installing-helm)
* [go](https://golang.org/doc)
* [dep](https://golang.github.io/dep/docs/introduction.html)
* [operator-sdk](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md)
* [minikube](https://github.com/kubernetes/minikube)

Ubuntu
```bash
# docker
sudo snap install docker

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

# minikube
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && \
  chmod +x minikube && \
  sudo mv minikube /usr/local/bin/
```

macOS
```bash
# docker
# download from https://hub.docker.com/editions/community/docker-ce-desktop-mac

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

# minikube
brew cask install minikube
```
