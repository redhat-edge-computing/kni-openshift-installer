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
	"os"
	"path"
	"path/filepath"
)

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "kni-install",
	RunE: rootCmdFunc,
}

// rootCmdFunc pre-configures system dependency locations prior to create() and destroy() calls
func rootCmdFunc(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type opts struct {
	isBareCluster bool
	isDryRun      bool
	kniRoot       string
	logLvl        string
	siteRepo      string
}

func (o opts) site() string {
	return path.Base(rootOpts.siteRepo)
}

func installer() string {
	return filepath.Join(rootOpts.kniRoot, rootOpts.site(), "requirements", "openshift-install")
}

func manifestDir() string {
	return filepath.Join(rootOpts.kniRoot, rootOpts.site(), "final_manifests")
}

var rootOpts = new(opts)

const flagSiteRepo = "site-repo"

func init() {

	userHome, _ := os.UserHomeDir()

	rootCmd.PersistentFlags().StringVar(&rootOpts.kniRoot, "kni-dir", filepath.Join(userHome, ".kni"), `(optional) Sets path to non-standard .kni path, useful for running the app outside of a containerized env.`)
	rootCmd.PersistentFlags().StringVar(&rootOpts.siteRepo, flagSiteRepo, "", `URI specifying path to site configs (e.g. github.com/path/to/site) (required)`)
	rootCmd.PersistentFlags().BoolVar(&rootOpts.isDryRun, "dry-run", false, `If true, prints but does not execute OS commands.`)
	rootCmd.PersistentFlags().StringVar(&rootOpts.logLvl, "log-level", "info", `Set log level of detail. Accepted input is one of: [debug | info | warn | error]`)
	rootCmd.PersistentFlags().BoolVar(&rootOpts.isBareCluster, "bare-cluster", false, "when true, complete cluster deployment and stop, do no deploy workload.")
}
