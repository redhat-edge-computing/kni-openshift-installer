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
	"net/url"
	"os/exec"
	"path/filepath"
	"k8s.io/klog"
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
		site, err := parseSite(blueprintUrl)
		if err != nil {
			panic(fmt.Errorf("failed to parse site from URL: %v", err))
		}
		klog.Infof("starting create for site %q", site)


		out, err := fetchRequirements(blueprintUrl)
		if err != nil {
			panic(fmt.Errorf("failed to fetch kni requirementes: %s\n%v", out, err))
		}
		klog.Infof("%v", out)

		out, err = prepareManifests(site)
		if err != nil {
			panic(fmt.Errorf("failed to prepare kni manifests: %s\n%v", out, err))
		}
		klog.Infof("%v", out)
	},
}

func fetchRequirements(blueprintUrl string) ([]byte, error) {
	klog.Infof("fetching site requirements from %q'", blueprintUrl)
	kniCmd := exec.Command("knictl", "fetch_requirements", blueprintUrl))
	return kniCmd.CombinedOutput()
}

func prepareManifests(site string) ([]byte, error) {
	klog.Infof("preparing manifests for site %q'", site)
	kniCmd := exec.Command("knictl", "prepare_manifests", site)
	return kniCmd.CombinedOutput()
}

func parseSite(siteLocation string) (string, error) {
	siteUrl, err := url.Parse(siteLocation)
	if err != nil {
		return "", fmt.Errorf("failed to extract path from URL: %v", err)
	}
	pathList := filepath.SplitList(siteUrl.Path)
	if len(pathList) == 0 {
		return "", fmt.Errorf("failed to extract site from path %q", siteUrl.Path)
	}
	return pathList[len(pathList) - 1], nil
}

func init() {
	rootCmd.AddCommand(createCmd)
}
