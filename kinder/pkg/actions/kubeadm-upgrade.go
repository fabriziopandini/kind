/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package actions

import (
	"fmt"
	"path/filepath"

	kcluster "sigs.k8s.io/kind/kinder/pkg/cluster"
)

//TODO: use consts for paths and filenames

// upgradeAction implements a developer friendly kubeadm upgrade workflow.
// please note that the upgrade will be executed by replacing kubeadm/kubelet/kubectl binaries;
// for sake of simplicity, we are skipping drain/uncordon when upgrading nodes
//
// this actions assumes that:
// 1) all the necessary images are already pre-loaded (otherwise kubeadm/kubelet will attempt to download images as usual)
// 2) the kubeadm/kubelet/kubectl binaries for the new kubernetes version are available in a well know place
//
// TODO:
// - apt upgrade, similar to user procedure (NB. currently only the apt mode uses deb during node-image creation, and the installation doesn't mark packages
// for preventing uncontrolled upgrades)
// - drain/uncordon of worker nodes
// - checking consistency of version among the provided binaries and the declared target version; if possible remove version flag
type upgradeAction struct{}

func init() {
	kcluster.RegisterAction("kubeadm-upgrade", newUpgradeAction)
}

func newUpgradeAction() kcluster.Action {
	return &upgradeAction{}
}

// Tasks returns the list of action tasks for the upgradeAction
func (b *upgradeAction) Tasks() []kcluster.Task {
	return []kcluster.Task{
		{
			Description: "Upgrade the kubeadm binary ⛵",
			TargetNodes: "@all",
			Run:         runUpgradeKubeadmBinary,
		},
		{
			Description: "Upgrade bootstrap control-plane ⛵",
			TargetNodes: "@cp1",
			Run:         runKubeadmUpgrade,
		},
		{
			Description: "Upgrade secondary control-planes ⛵",
			TargetNodes: "@cpN",
			Run:         runKubeadmUpgradeControlPlane,
		},
		{
			Description: "Upgrade workers config ⛵",
			TargetNodes: "@w*",
			Run:         runKubeadmUpgradeWorkers,
		},
		{
			Description: "Upgrade kubelet and kubectl ⛵",
			TargetNodes: "@all",
			Run:         runUpgradeKubeletKubectl,
		},
	}
}

func runUpgradeKubeadmBinary(kctx *kcluster.KContext, kn *kcluster.KNode, flags kcluster.ActionFlags) error {

	src := filepath.Join("/kinder", "upgrade", "kubeadm")
	dest := filepath.Join("/usr", "bin", "kubeadm")

	fmt.Println("==> upgrading kubeadm 🚀")
	if err := kn.Command(
		"cp", src, dest,
	).Run(); err != nil {
		return err
	}

	return nil
}

func runKubeadmUpgrade(kctx *kcluster.KContext, kn *kcluster.KNode, flags kcluster.ActionFlags) error {
	if err := kn.DebugCmd(
		"==> kubeadm upgrade apply 🚀",
		"kubeadm", "upgrade", "apply", "-f", flags.UpgradeVersion.String(),
	); err != nil {
		return err
	}

	//TODO: check if download config included (and if restart kubelet included)

	return nil
}

func runKubeadmUpgradeControlPlane(kctx *kcluster.KContext, kn *kcluster.KNode, flags kcluster.ActionFlags) error {
	if err := kn.DebugCmd(
		"==> kubeadm upgrade node experimental-control-plane 🚀",
		"kubeadm", "upgrade", "node", "experimental-control-plane",
	); err != nil {
		return err
	}

	//TODO: check if download config included (and if restart kubelet included)

	return nil
}

func runKubeadmUpgradeWorkers(kctx *kcluster.KContext, kn *kcluster.KNode, flags kcluster.ActionFlags) error {
	if err := kn.DebugCmd(
		"==> kubeadm upgrade node config 🚀",
		"kubeadm", "upgrade", "node", "config", "--kubelet-version", flags.UpgradeVersion.String(),
	); err != nil {
		return err
	}

	//TODO: check if restart kubelet included

	return nil
}

func runUpgradeKubeletKubectl(kctx *kcluster.KContext, kn *kcluster.KNode, flags kcluster.ActionFlags) error {
	// upgrade kubectl
	fmt.Println("==> upgrading kubectl 🚀")
	src := filepath.Join("/kinder", "upgrade", "kubectl")
	dest := filepath.Join("/usr", "bin", "kubectl")

	if err := kn.Command(
		"cp", src, dest,
	).Run(); err != nil {
		return err
	}

	// upgrade kubelet
	fmt.Println("==> upgrading kubelet 🚀")
	src = filepath.Join("/kinder", "upgrade", "kubelet")
	dest = filepath.Join("/usr", "bin", "kubelet")

	if err := kn.Command(
		"cp", src, dest,
	).Run(); err != nil {
		return err
	}

	fmt.Println("==> restart kubelet 🚀")
	if err := kn.Command(
		"systemctl", "restart", "kubelet",
	).Run(); err != nil {
		return err
	}

	//write "/kind/version"
	if err := kn.Command(
		"echo", fmt.Sprintf("\"%s\"", flags.UpgradeVersion.String()), ">", "/kind/version",
	).Run(); err != nil {
		return err
	}

	return nil
}
