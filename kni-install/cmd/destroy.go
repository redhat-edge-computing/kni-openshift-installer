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
	"github.com/spf13/cobra"
	"k8s.io/klog"
	"os/exec"
)

// destroyCmd represents the destroy command
var (
	destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	destroyClusterCmd = &cobra.Command{
		Use:   "cluster",
		Short: "A brief description of your command",
		RunE:   destroyCluster,
	}
)

func destroyCluster(_ *cobra.Command, _ []string) error {
	klog.Info("tearing down site cluster")
	err := execCmdToStdout(exec.Command(ocpInstaller, "destroy", "cluster", "--log-level=debug", "--dir", siteBuildDir))
	if err != nil {
		return fmt.Errorf("destroy cluster failed: %v", err)
	}
	return nil
}

func init() {
	destroyCmd.AddCommand(destroyClusterCmd)
	rootCmd.AddCommand(destroyCmd)
}
