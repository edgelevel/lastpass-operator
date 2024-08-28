# Setup

* [docker](https://docs.docker.com/install)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)
* [helm](https://helm.sh/docs/using_helm/#installing-helm)
* [go](https://golang.org/doc)
* [operator-sdk](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md)
* [minikube](https://github.com/kubernetes/minikube)

## Ubuntu:

Install docker using the instructions below from the [docker website](https://docs.docker.com/engine/install/ubuntu/#installation-methods)
```bash
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

### kubectl: 

`sudo snap install kubectl --classic`

### helm: 

`sudo snap install helm --classic`

### go:

`sudo snap install go --classic`


### Operator-SDK:
```bash
# operator-sdk (from source)
go get -d github.com/operator-framework/operator-sdk
cd $GOPATH/src/github.com/operator-framework/operator-sdk
git checkout master
make dep
make install
```
Alternatively, you can use the operator-sdk build dependency that is in the `bin` directory. This operator-sdk binary can be downloaded using:

`make operator-sdk`

### minikube
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64
```

## macOS
```bash
# docker
# download from https://hub.docker.com/editions/community/docker-ce-desktop-mac

# kubectl
brew install kubernetes-cli

# helm
brew install kubernetes-helm

# go
brew install go

# operator-sdk
brew install operator-sdk

# minikube
brew cask install minikube
```
