/*
Copyright © 2020 Jonathan Cope jcope@redhat.com

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
package cmd

import (
	"fmt"
	"os/exec"

	"k8s.io/klog"

	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var (
	destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "this command does not do anything on its own.  you must call the cluster" +
			"subcommand in order to target a cluster for tear down",
	}

	destroyClusterCmd = &cobra.Command{
		Use:   "cluster",
		Short: "tear down the cluster associated with the given site",
		Long: "wraps openshift-install to completely tear down the cluster" +
			"deployed from the blueprint for the given site",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return verifyRequiredFlags(cmd)
		},
		RunE:  destroyCluster,
	}
)

func destroyCluster(cmd *cobra.Command, _ []string) error {

	klog.Info("tearing down site cluster")
	err := execCmdToStdout(exec.Command(installer(), "destroy", "cluster", "--log-level=debug", "--dir", manifestDir()))
	if err != nil {
		return fmt.Errorf("destroy cluster failed: %v", err)
	}
	return nil
}

func init() {
	destroyCmd.AddCommand(destroyClusterCmd)
	rootCmd.AddCommand(destroyCmd)
}
