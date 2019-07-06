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

# build and publish
# TODO https://hub.docker.com/search/?q=niqdev&type=image
make docker-push tag=0.1.0

# apply chart
kubectl create namespace lastpass
helm template \
  --values chart/values.yaml \
  --set lastpass.username="myUsername" \
  --set lastpass.password="myPassword" \
  chart/ | kubectl apply -n lastpass -f -
```
