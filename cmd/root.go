
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gokmp",
	Short: "gokmp: Backup your Flickr pictures",
	Long: `gokmp: Automation tool for backing up your Flickr pictures locally on disk`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}

func init() {
	// nil
}
