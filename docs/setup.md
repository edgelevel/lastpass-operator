## Setup

* [docker](https://docs.docker.com/install)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)
* [helm](https://helm.sh/docs/intro/install)
* [go](https://go.dev/doc/install)
* [operator-sdk](https://sdk.operatorframework.io/docs/installation)
* [minikube](https://minikube.sigs.k8s.io/docs)

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

# operator-sdk (from source)
git clone https://github.com/operator-framework/operator-sdk
cd operator-sdk
git checkout master
make install

# minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb && \
  sudo dpkg -i minikube_latest_amd64.deb
```

macOS
```bash
# docker
# download from https://docs.docker.com/desktop/mac/install

# kubectl
brew install kubernetes-cli

# helm
brew install kubernetes-helm

# go
brew install go

# operator-sdk
brew install operator-sdk

# minikube
brew install minikube
```
