# lastpass-operator

A Kubernetes Operator to manage [secrets](https://kubernetes.io/docs/concepts/configuration/secret) stored in [LastPass](https://www.lastpass.com) password manager

## How it works

Suppose you have some credentials stored in LastPass
```bash
$ lpass show example/my-secret --json
[
  {
    "id": "8190226423897406876",
    "name": "my-secret",
    "fullname": "example/my-secret",
    "username": "whoami",
    "password": "s3cr3t",
    "last_modified_gmt": "1562690587",
    "last_touch": "0",
    "group": "example",
    "url": "https://lastpass.com",
    "note": "{\"myKey\":\"myValue\"}"
  }
]
```

Define a `LastPass` [Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources) to automatically manage the lifecycle of your secrets in Kubernetes
```bash
$ cat example/niqdev_v1alpha1_lastpass_cr.yaml
apiVersion: niqdev.com/v1alpha1
kind: LastPass
metadata:
  name: example-lastpass
spec:
  secretRef:
    group: example
    name: my-secret
    withUsername: true
    withPassword: true
    withUrl: true
    withNote: true
  # TODO
  syncPolicy:
    enable: true
    refresh: 60

# create custom resource
$ kubectl apply -f example/niqdev_v1alpha1_lastpass_cr.yaml
```

The operator will take care of create native Kubernetes secrets and keep them up to date that if they change
```bash
# verify
$ kubectl get lastpass
$ kubectl get secrets

# inspect
$ kubectl get secret example-lastpass-8190226423897406876 -o yaml
apiVersion: v1
data:
  NOTE: eyJteUtleSI6Im15VmFsdWUifQ==
  PASSWORD: czNjcjN0
  URL: aHR0cHM6Ly9sYXN0cGFzcy5jb20=
  USERNAME: d2hvYW1p
kind: Secret
metadata:
  annotations:
    fullname: example/my-secret
    group: example
    id: "8190226423897406876"
    lastModifiedGmt: "1562690587"
    lastTouch: "0"
    name: my-secret
  creationTimestamp: "2019-07-09T15:00:13Z"
  labels:
    app: lastpass-operator
  name: example-lastpass-8190226423897406876
  namespace: default
  ownerReferences:
  - apiVersion: niqdev.com/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: LastPass
    name: example-lastpass
    uid: 0687d5a7-5f02-4ee4-a6c4-011c734f4149
  resourceVersion: "113312"
  selfLink: /api/v1/namespaces/default/secrets/example-lastpass-8190226423897406876
  uid: 382008d2-8999-444d-86c8-e4f29eecbe9f
type: Opaque

# check values
$ echo 'czNjcjN0' | base64 --decode
s3cr3t
$ echo 'eyJteUtleSI6Im15VmFsdWUifQ==' | base64 --decode | jq -c
{"myKey":"myValue"}
```

## Development

* [Setup](doc/setup.md)
* [golang](doc/golang.md)
* [operator-sdk](doc/operator.md)
* [LastPass](doc/lastpass.md)

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

TODO
* [ ] fix `lpass` permissions in Dockerfile
* [ ] fix rules in `rbac.yaml`
* [ ] version of lastpass-cli
