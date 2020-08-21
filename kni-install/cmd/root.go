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
	"github.com/spf13/viper"
	"k8s.io/klog"
	"os"
	"path"
	"path/filepath"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kni-install",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	isDryRun     bool
	isBareCluster  bool
	logLvl       string
	site         string
	siteRepo     string
	siteBuildDir string
	ocpInstaller string
)

func init() {
	cobra.OnInitialize(initConfig)

	userHome, _ := os.UserHomeDir()

	kniRoot := rootCmd.PersistentFlags().String("kni-dir", filepath.Join(userHome, ".kni"), `(optional) Sets path to non-standard .kni path`)
	rootCmd.PersistentFlags().StringVar(&siteRepo, "repo", "", `git repo path containing site config files`)
	rootCmd.PersistentFlags().BoolVar(&isDryRun, "dry-run", false, `(optional) If true, prints, but does not execute OS commands.`)
	rootCmd.PersistentFlags().StringVar(&logLvl, "log-level", "info", `Set log level of detail. Accepted input is one of: ["info", "debug"]`)
	rootCmd.PersistentFlags().BoolVar(&isBareCluster, "bare-cluster", false, "when true, complete cluster deployment and stop, do no deploy workload.")
	_ = rootCmd.PersistentFlags().Parse(os.Args[1:])

	_, err := os.Stat(*kniRoot)
	if err != nil {
		klog.Fatalf("stat failed for dir %q: %v", *kniRoot, err)
	}
	site = path.Base(siteRepo)
	siteBuildDir = filepath.Join(*kniRoot, site, "final_manifests")
	ocpInstaller = filepath.Join(*kniRoot, site, "requirements", "openshift-install")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kni-openshift-installer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kni-openshift-kni-install")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
