## operator-sdk

* [operator-sdk](https://sdk.operatorframework.io/docs/building-operators/golang)
* [Controller Runtime Client API](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/client.md)
* [A complete guide to Kubernetes Operator SDK](https://banzaicloud.com/blog/operator-sdk)
* [OperatorHub](https://operatorhub.io)
* [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions)
* [Groups and Versions and Kinds, oh my!](https://book.kubebuilder.io/cronjob-tutorial/gvks.html)
* [Level Triggering and Reconciliation in Kubernetes](https://hackernoon.com/level-triggering-and-reconciliation-in-kubernetes-1f17fe30333d)
* [A deep dive into Kubernetes controllers](https://engineering.bitnami.com/articles/a-deep-dive-into-kubernetes-controllers.html)

Initialize project
```bash
# create project
mkdir -p $GOPATH/src/github.com/edgelevel && cd $_
operator-sdk new lastpass-operator --dep-manager=dep
operator-sdk init --domain example.com --repo github.com/example/memcached-operator

# add crd
operator-sdk add api --api-version=edgelevel.com/v1alpha1 --kind=LastPass
operator-sdk generate k8s

# add controller
operator-sdk add controller --api-version=edgelevel.com/v1alpha1 --kind=LastPass
```
