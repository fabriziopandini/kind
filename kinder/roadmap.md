# Roadmap 🗺️

This document outlines some goals, non-goals, and future aspirations for kinder.

High level goals for kinder include:

- [ ] Enable usage of kinder(kind) as local integration test environment for kubeadm/cluster api development
   - [x] Allow creation of nodes "ready for installing Kubernetes"
   - [x] Provide pre built “developer” workflows for kubedam init, join, reset
      - [x] init and init with phases
      - [x] join and join with phases
      - [x] reset
      - [x] init and join with automatic copy certs 
      - [x] Provide pre built “developer” workflow for kubeadm upgrades
   - [x] Allow build node-image variants
      - [x] add pre loaded images
      - [x] replace the kubeadm binary
      - [x] add kubernetes binaries for upgrades
   - [x] Allow test of kubeadm cluster variations 
      - [x] external etcd
      - [x] kube-dns
   - [x] Provide "topology aware" wrappers for `docker exec` and `docker cp`
   - [ ] Provide a way to add nodes to an existing cluster
      - [x] Add worker node
      - [ ] Add control plane node (and reconfigure load balancer) 
   - [ ] Provide smoke test action
   - [ ] E2E run actions

**Non**-Goals include:

- Replace or fork kind. kind is awesome and we are committed to help to improve it (see long term goals)
- Supporting every possible use cases that can be build on top of kind as a library

Longer Term goals include:

- Simplify daily activities for kubeadm/cluster api development
- Help new contributors starting to work on kubeadm/cluster api development
- Contribute to improving and testing "kind as a library"
- Contribute idea/code for new features in kind
- Provide a home for use cases that are difficult to support in the main kind CLI
