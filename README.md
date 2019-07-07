# lastpass-operator

> TODO

* [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret)

## Development

* [LastPass](doc/lastpass.md)
* [operator-sdk](doc/operator.md)

```bash
# download source
mkdir -p $GOPATH/src/github.com/niqdev && cd $_
git clone git@github.com:niqdev/lastpass-operator.git
cd lastpass-operator

# install dependencies
dep ensure
```

Run locally outside the cluster on [minkube](https://github.com/kubernetes/minikube)
```bash
# requires virtualbox
minikube start

# apply crd
kubectl apply -f chart/templates/crd.yaml

# run locally
export OPERATOR_NAME=lastpass-operator
operator-sdk up local --namespace=default --verbose
```

Run as a Deployment inside the cluster
```bash
# build and publish
# https://hub.docker.com/repository/docker/niqdev/lastpass-operator
make docker-push tag=0.1.0

# apply chart
kubectl create namespace lastpass
helm template \
  --values chart/values.yaml \
  --set lastpass.username="myUsername" \
  --set lastpass.password="myPassword" \
  chart/ | kubectl apply -n lastpass -f -
```

TODO
* [ ] fix `lpass` permissions in Dockerfile
* [ ] fix rules in `rbac.yaml`
