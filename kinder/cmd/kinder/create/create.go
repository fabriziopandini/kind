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

// Package create implements the `create` command
// Nb. re-implemented in Kinder in order to import the kinder version of createcluster
package create

import (
	"github.com/spf13/cobra"

	createcluster "sigs.k8s.io/kind/kinder/cmd/kinder/create/cluster"
	createnode "sigs.k8s.io/kind/kinder/cmd/kinder/create/node"
	"sigs.k8s.io/kind/pkg/cluster/config"
)

// NewCommand returns a new cobra.Command for cluster creation
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "Creates one of [cluster, worker-node, control-plane-node]",
		Long:  "Creates one of local Kubernetes cluster (cluster), or nodes in a local kubernetes cluster (worker-node, control-plane-node)",
	}
	cmd.AddCommand(createcluster.NewCommand())
	cmd.AddCommand(createnode.NewCommand(config.ControlPlaneRole)) //this is currently super hacky, looking for a better solution
	cmd.AddCommand(createnode.NewCommand(config.WorkerRole))
	return cmd
}
