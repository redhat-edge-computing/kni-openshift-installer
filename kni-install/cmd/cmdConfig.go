package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
)

func subCmdPreConfig(_ *cobra.Command, _ []string) error {
	if len(rootOpts.siteRepo) == 0 {
		return fmt.Errorf("required flag %q no set", flagSiteRepo)
	}
	return nil
}

func site() string {
	return path.Base(rootOpts.siteRepo)
}

func installer() string {
	return filepath.Join(rootOpts.kniRoot, site(), "requirements", "openshift-install")
}

func manifestDir() string {
	return filepath.Join(rootOpts.kniRoot, site(), "final_manifests")
}
