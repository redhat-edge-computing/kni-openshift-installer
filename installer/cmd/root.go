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
	"net/url"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "installer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const FLAG_URL = "url"
const FLAG_DRYRUN = "dry-run"

const kniRoot = "/root/.kni"

var (
	siteDomain       string
	requirementsPath string
	installerPath    string
	manifestsPath    string
	isDryRun         *bool // TODO implement dry run
)

var _ = isDryRun // short circuit unused var error until I get that implemented

func init() {
	cobra.OnInitialize(initConfig)

	sourceURL := rootCmd.PersistentFlags().String(FLAG_URL, "", "the URL (http, file, etc) containing the site blueprints")
	isDryRun = rootCmd.PersistentFlags().Bool(FLAG_DRYRUN, false, "if true, only print the os commands")
	if err := rootCmd.PersistentFlags().Parse(os.Args[1:]); err != nil {
		klog.Fatalf("flags error: %v", err)
	}

	var err error
	siteDomain, err = parseSite(*sourceURL)
	if err != nil {
		klog.Fatalf("could not parse site from url: %v", err)
	}

	requirementsPath = filepath.Join(kniRoot, siteDomain, "requirements")
	installerPath = filepath.Join(requirementsPath, "openshift-install")
	manifestsPath = filepath.Join(kniRoot, siteDomain, "final_manifests")
}

func parseSite(siteLocation string) (string, error) {
	siteUrl, err := url.Parse(siteLocation)
	if err != nil {
		return "", fmt.Errorf("failed to extract path from URL: %v", err)
	}
	// trailing '/' break path splitting, trim them to be safe
	siteUrl.Path = strings.Trim(siteUrl.Path, "/")
	pathList := strings.Split(siteUrl.Path, "/")
	if len(pathList) == 0 {
		return "", fmt.Errorf("failed to extract site from path %q", siteUrl.Path)
	}
	return pathList[len(pathList)-1], nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".installer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".installer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
