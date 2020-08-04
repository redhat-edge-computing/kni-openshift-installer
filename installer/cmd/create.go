/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	createClusterCmd = &cobra.Command{
		Use: "cluster",
		Short: "",
		Long: "",
		Run: createCluster,
		Args: cobra.ExactArgs(0),
	}

	createIgnitionConfigsCmd = &cobra.Command{
		Use:   "ignition-config",
		Short: "",
		Long:  "",
		Run:   createIgnitionConfigs,
		Args:  cobra.ExactArgs(0),
	}
)

func init() {
	rootCmd.AddCommand(newCreateCmd())
}

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "A brief description of your command",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(createClusterCmd)
	cmd.AddCommand(createIgnitionConfigsCmd)
	return cmd
}

func createCluster(cmd *cobra.Command, args []string) {
		blueprintUrl := cmd.PersistentFlags().Lookup(FLAG_URL).Value.String()
		site := siteDomain // to avoid numerous references to a global var
		klog.Infof("starting create for site %q", site)

		err := fetchRequirements(blueprintUrl)
		if err != nil {
			klog.Fatalf("failed to fetch kni requirements: %v", err)
		}

		err = prepareManifests(site)
		if err != nil {
			klog.Fatalf("failed to prepare kni manifests: %v", err)
		}

		err = deployCluster(site)
		if err != nil {
			klog.Fatalf("failed to deploy cluster: %v", err)
		}
}

func fetchRequirements(blueprintUrl string) error {
	klog.Infof("fetching site requirements")
	return execCmdToStdout(exec.Command("knictl", "fetch_requirements", blueprintUrl))
}

func prepareManifests(site string) error {
	klog.Info("preparing manifests")
	return execCmdToStdout(exec.Command("knictl", "prepare_manifests", site))
}

func deployCluster(site string) error {
	klog.Infof("deploying cluster")
	return execCmdToStdout(exec.Command(installerPath, "create", "cluster", "--log-level=debug", "--dir", manifestsPath))
}

func execCmdToStdout(command *exec.Cmd) error {
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Start()
	if err != nil {
		return fmt.Errorf("exec failed: %v", err)
	}
	return command.Wait()
}

// stub
func createIgnitionConfigs(cmd *cobra.Command, args []string) {
	return
}