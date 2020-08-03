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
	"github.com/spf13/cobra"
	"k8s.io/klog"
	"os"
	"os/exec"
)

// docker run --rm  -v $HOME/.aws:/root/.aws:Z quay.io/jcope/cluster create -site $SITE

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {

		blueprintUrl := rootCmd.PersistentFlags().Lookup(FLAG_URL).Value.String()
		site := siteDomain // to avoid numerous references to a global var
		klog.Infof("starting create for site %q", site)

		err := fetchRequirements(blueprintUrl)
		if err != nil {
			klog.Fatalf("failed to fetch kni requirementes: %v", err)
		}

		err = prepareManifests(site)
		if err != nil {
			klog.Fatalf("failed to prepare kni manifests: %v", err)
		}

		err = deployCluster(site)
		if err != nil {
			klog.Fatalf("failed to deploy cluster: %v", err)
		}
	},
}

func fetchRequirements(blueprintUrl string) error {
	klog.Infof("fetching site requirements")
	kniCmd := exec.Command("knictl", "fetch_requirements", blueprintUrl)
	klog.Infof("exec: %v", kniCmd.String())
	kniCmd.Stderr = os.Stderr
	kniCmd.Stdout = os.Stdout
	err := kniCmd.Start()
	if err != nil {
		return err
	}
	return kniCmd.Wait()
}

func prepareManifests(site string) error {
	klog.Info("preparing manifests")
	kniCmd := exec.Command("knictl", "prepare_manifests", site)
	klog.Infof("exec: %v", kniCmd.String())
	kniCmd.Stderr = os.Stderr
	kniCmd.Stdout = os.Stdout
	err := kniCmd.Start()
	if err != nil {
		return err
	}
	return kniCmd.Wait()
}

func deployCluster(site string) error {
	klog.Infof("deploying cluster")
	ocpInstallCmd := exec.Command(installerPath, "create", "cluster", "--log-level=debug", "--dir", manifestsPath)
	klog.Infof("exec: %v", ocpInstallCmd.String())
	ocpInstallCmd.Stdout = os.Stdout
	ocpInstallCmd.Stderr = os.Stderr
	err := ocpInstallCmd.Start()
	if err != nil {
		klog.Fatalf("failed to exec openshift create cluster %s: %v", siteDomain, err)
	}
	return ocpInstallCmd.Wait()
}

func init() {
	rootCmd.AddCommand(createCmd)
}
