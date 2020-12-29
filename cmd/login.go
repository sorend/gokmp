package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sorend/gokmp/pkg/login"
)

var loginCmd = &cobra.Command{
	Use: "login",
	Short: "Login to Flickr",
	Long: `Login to Flickr with 'read' rights in order to backup.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return login.Run(FlickrApiKey, FlickrApiSecret)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
