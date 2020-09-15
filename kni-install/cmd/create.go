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
	"os"
	"os/exec"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRunE: subCmdPreConfig,
	}

	createClusterCmd = &cobra.Command{
		Use:   "cluster",
		Short: "deploys a cluster per the specified site",
		Long: `Wraps multiple knictl command line executions to fetch requirements,
prepare manifests, create the cluster, and apply workloads as defined in the specified site.`,
		RunE:    execCreateCmd,
		Args:    cobra.ExactArgs(0),
		PreRunE: initCreateCommand,
	}

	createIgnitionConfigsCmd = &cobra.Command{
		Use:   "ignition-configs",
		Short: "prepares ignition config manifests for baremetal deployments. does not create a cluster",
		Long: `Wraps multiple knictl command line executions to fetch requirements,
prepare manifests, create the ignition-configs to be used for baremetal deployments.`,
		RunE:     execCreateIgnitionConfigsCmd,
		Args:    cobra.ExactArgs(0),
		PreRunE: initCreateCommand,
	}
)

func initCreateCommand(_ *cobra.Command, _ []string) error {
	err := fetchRequirements()
	if err != nil {
		return err
	}

	err = prepareManifests()
	if err != nil {
		return err
	}
	return nil
}

func fetchRequirements() error {
	klog.Info("fetching site requirements")
	err := execCmdToStdout(exec.Command("knictl", "fetch_requirements", rootOpts.siteRepo))
	if err != nil {
		return fmt.Errorf("fetching requirements failed: %v", err)
	}
	klog.Info("done fetching requirements")
	return nil
}

func prepareManifests() error {
	klog.Info("preparing manifests")
	err := execCmdToStdout(exec.Command("knictl", "prepare_manifests", site()))
	if err != nil {
		return fmt.Errorf("manifest preparation failed: %v", err)
	}
	klog.Info("done preparing manifests")
	return nil
}

func execCreateCmd(_ *cobra.Command, _ []string) error {
	err := createCluster()
	if err != nil {
		return err
	}

	if ! rootOpts.isBareCluster {
		err = applyWorkloads()
		if err != nil {
			return err
		}
	}
	return nil
}

func createCluster() (err error) {
	klog.Info("deploy cluster")
	err = execCmdToStdout(exec.Command(installer(), "create", "cluster", "--log-level", rootOpts.logLvl, "--dir", manifestDir()))
	if err != nil {
		return fmt.Errorf("cluster deployment failed: %v", err)
	}
	klog.Info("cluster deployment complete")
	return nil
}

func applyWorkloads() (err error) {
	klog.Info("applying workload manifests")
	err = execCmdToStdout(exec.Command("knictl", "apply_workloads", site()))
	if err != nil {
		return fmt.Errorf("apply workloads failed: %s", err)
	}
	klog.Info("workload manifests deployed")
	return nil
}

func execCreateIgnitionConfigsCmd(_ *cobra.Command, _ []string) error {
	klog.Info("creating ignition-configs")
	err := execCmdToStdout(exec.Command(installer(), "create", "ignition-configs", "--log-level", rootOpts.logLvl, "--dir", manifestDir()))
	if err != nil {
		return fmt.Errorf("create ignition configs failed: %v", err)
	}
	klog.Info("ignition-configs creation complete")
	return nil
}

func execCmdToStdout(command *exec.Cmd) error {
	if rootOpts.isDryRun {
		klog.Infof("dry-run exec: %s", command.String())
		return nil
	}
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Start()
	if err != nil {
		return fmt.Errorf("command failed: %v", err)
	}
	return command.Wait()
}

func init() {
	createCmd.AddCommand(createClusterCmd)
	createCmd.AddCommand(createIgnitionConfigsCmd)
	rootCmd.AddCommand(createCmd)
}
