## operator-sdk

* [operator-sdk](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md)
* [A complete guide to Kubernetes Operator SDK](https://banzaicloud.com/blog/operator-sdk)
* [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions)
* [Groups and Versions and Kinds, oh my!](https://book.kubebuilder.io/cronjob-tutorial/gvks.html)
* [Controller Runtime Client API](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/client.md)
* [OperatorHub](https://operatorhub.io)

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
