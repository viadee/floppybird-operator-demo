# floppybird-operator

## prerequisite

- go version v1.19.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- k3d (or any other local k8s cluster solution)

## local test

```shell
# create local cluster
k3d cluster create local-cluster
kubectl config use-context k3d-local-cluster

# install CRD
make install

## launch from vscode "Run and Debug"
```

## project initialization

```shell
# create project
kubebuilder init --domain demo.viadee.de --repo github.com/viadee/floppybird-operator-demo
# create api (Create Resource [y/n] y; Create Controller [y/n] y)
kubebuilder create api --group webapp --version v1alpha1 --kind Floppybird
```
