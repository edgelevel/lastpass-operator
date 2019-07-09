# lastpass-operator

> TODO

* TODO version of lastpass-cli

* [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret)

## Development

* [LastPass](doc/lastpass.md)
* [Setup](doc/setup.md)
* [golang](doc/golang.md)
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
export LASTPASS_USERNAME=myUsername
export LASTPASS_PASSWORD=myPassword
operator-sdk up local --namespace=default --verbose
```

Run as a Deployment inside the cluster
```bash
# build and publish
# https://hub.docker.com/repository/docker/niqdev/lastpass-operator
make docker-push tag=0.2.0

# apply chart
kubectl create namespace lastpass
helm template \
  --values chart/values.yaml \
  --set lastpass.username="myUsername" \
  --set lastpass.password="myPassword" \
  chart/ | kubectl apply -n lastpass -f -
```

---

TODO

```bash
kubectl apply -f example/niqdev_v1alpha1_lastpass_cr.yaml

kubectl get secrets
kubectl get secret example-lastpass-<SECRET_ID> -o yaml
echo '' | base64 --decode
```

TODO
* [ ] fix `lpass` permissions in Dockerfile
* [ ] fix rules in `rbac.yaml`
