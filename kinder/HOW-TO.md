# How to

This document shows how to use kinder for kubeadm manual testing.

## Preparing for tests

Before starting test, it necessary to get a node-image to be used as a base for nodes in the cluster.

There are three ways to get an image:

## Use a kind public node-images

The esiest way to get a node image for a major/minor Kubernetes version is use kind images available on docker hub, e.g.

```bash
docker pull kindest/node:vX
```

## build a node-image

For building a node image you can refer to kind documentation; below a short recap of necessary steps:

Build a base-image (or download one from docker hub)
```bash
kinder build base-image --image kindest/base:latest
```

Build a node-image starting from the above node image
```bash
# To build a node-image using latest Kubernetes apt packages available
kinder build node-image --base-image kindest/base:latest --image kindest/node:vX --type apt

# To build a node-image using local Kubernetes repository
kinder build node-image --base-image kindest/base:latest --image kindest/node:vX --type bazel
```

>  NB see https://github.com/kubernetes/kubeadm/blob/master/testing-pre-releases.md#change-the-target-version-number-when-building-a-local-release for overriding the build version in case of `--type bazel`

## customize a node-image

As a third option for building node-image, it is possible to pick an existing node image and customize it by:

1. overriding the kubeadm/kubelet/kubectl binary, and eventually the `/kind/version` file

TODO: example

2. adding/overriding the pre loaded images in the `/kind/images` folder

TODO: example

3. adding a second kubernetes version in `...` for testing upgrades

TODO: example

## Test kubeadm version vX

The document assumes vX the kubeadm release to be tested.

Test includes verification of kubeadm init and join workflows and supported cluster variants.

The test can be executed using a node-image vX.

### Test single master cluster

```bash
# create a cluster with eventually some worker node
kinder create cluster --image kindest/node:vX --workers-nodes 0

# initialize the bootstrap control plane
kinder do kubeadm-init

# join secondary control planes (if any)
kinder do kubeadm-join
```

test variants:

1. add `--kube-dns` flag to `kinder create cluster` to test usage of kube-dns instead of CoreDNS
2. add `--external-etcd` flag to `kinder create cluster` to test usage of external etcd cluster
3. add `--use-phases` flag to `kubeadm-init` and/or `kubeadm-join` to test phases
4. any combination of the above

validation and cleanup

```bash
# verify kubeadm commands outputs

# verify the resulting cluster
kinder do cluster-info
# > check for nodes, Kubernetes version x, ready
# > check all the components running, Kubernetes version x + related dependencies
# > check for etcd member

# cleanup
kinder do kubeadm-reset

kinder delete cluster
```

### Test HA cluster

```bash
# create a cluster with at least two control plane nodes (and eventually some worker node)
kinder create cluster --image kindest/node:vX --control-plane-nodes 2 --workers-nodes 0

# initialize the bootstrap control plane
kinder do kubeadm-init

# join secondary control planes
kinder do kubeadm-join
```

test variants:

1. add `--kube-dns` flag to `kinder create cluster` to test usage of kube-dns instead of CoreDNS
2. add `--external-etcd` flag to `kinder create cluster` to test usage of external etcd cluster
3. add `--use-phases` flag to `kubeadm-init` and/or `kubeadm-join` to test phases
4. add `--automatic-copy-certs` flag both to `kubeadm-init` and `kubeadm-join` to test the automatic copy certs feature
4. any combination of the above

validation and cleanup

```bash
# verify kubeadm commands outputs

# verify the resulting cluster
kinder do cluster-info
# > check for nodes, Kubernetes version x, ready
# > check all the components running, Kubernetes version x + related dependencies
# > check for etcd members

# cleanup
kinder do kubeadm-reset

kinder delete cluster
```

## Test upgrades from vX-1 to vX

kubeadm vX is espected to be capable to upgrade Kubernetes cluster version vX-1.

The test can be executed using a kind node-image vX-1 with a second kubernetes version vX deployed in the `/kinder/upgrade`; this can be achieved with the following command

```bash
# Add vX+1 binaries and images
kinder build node-variant --base-image kindest/node:vX-1 --image kindest/node:vX-1toX --with-upgrade-binaries $working_dir/packages/vX/binaries --with-images $working_dir/packages/vX/images
```

Once the image is ready, same test/test variants of Test kubeadm version vX apply for the initial part of the test (from vX-1). 


After it is possible to execute upgrade (to vX)


```bash
# upgrade the cluster
kinder do kubeadm-upgrade

# verify the resulting cluster
kinder do cluster-info
# > check for nodes, Kubernetes version x, ready
# > check all the components running, Kubernetes version x + related dependencies
# > check for etcd members
```

Eventually it is possible to add/join additional nodes and check the state of the resulting cluster

TODO: example

## Test kubeadm version vX on Kubernetes version vX-1

kubeadm vX is espected to be capable to install and manage Kubernetes version vX-1.

The test can be executed customizing a kind node-image vX-1 with a kubeadm binary vX.

TODO: example

Once the image is ready, same test/test variants of Test kubeadm version vX apply

## Test upgrades from vX to vX+1

Before vX release cut it is important to test upgrades to future minor release, trying yo anticipate potential problems
and consequent fix on the vX branch.

The test can be executed using a kind node-image vX with a second kubernetes version vX+1 in `...` for testing upgrades
kubernetes version vX+1 (that does not exists yet) can be created as a build of vX with a different build tag.

```bash
# Add vX+1 binaries
kinder build node-variant --with-upgrade-binaries vX+1/binaries

# Add vX+1 images
kinder build node-variant --with-images vX+1/images
```

Once the image is ready, same test/test variants of Test kubeadm version vX apply for the initial part of the test (from vX). After it is possible to execute upgrade (to vX+1):


```bash
# upgrade the cluster
kinder do kubeadm-upgrade

# verify the resulting cluster
kinder do cluster-info
# > check for nodes, Kubernetes version x, ready
# > check all the components running, Kubernetes version x + related dependencies
# > check for etcd members
```

Eventually it is possible to add/join additional nodes and check the state of the resulting cluster.

TODO: example