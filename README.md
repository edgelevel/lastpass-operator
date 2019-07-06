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

# build operator
dep ensure
make docker-build tag=<VERSION>

# TODO login + publish
# https://hub.docker.com/search/?q=niqdev&type=image

# apply chart
kubectl create namespace lastpass
helm template \
  --values chart/values.yaml \
  --set lastpass.username="myUsername" \
  --set lastpass.password="myPassword" \
  chart/ | kubectl apply -n lastpass -f -
```
