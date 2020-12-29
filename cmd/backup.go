package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sorend/gokmp/pkg/backup"
)

var BackupDirectory string
var BackupNsid string

var backupCmd = &cobra.Command{
	Use: "backup",
	Short: "Backup from Flickr",
	Long: `Backup pictures from Flickr to local disk`,
	RunE: func(cmd *cobra.Command, args []string) error {
		accessToken := viper.GetString("accessToken")
		accessSecret := viper.GetString("accessSecret")
		if accessToken == "" || accessSecret == "" {
			return errors.New("Please 'login' before calling 'backup' (no access token found)")
		}
		if err := backup.Run(FlickrApiKey, FlickrApiSecret, accessToken, accessSecret, BackupNsid, BackupDirectory); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	backupCmd.Flags().StringVarP(&BackupDirectory, "destination", "d", "", "Destination for backup")
	backupCmd.Flags().StringVarP(&BackupNsid, "nsid", "i", "", "Flickr user@id to backup")
	rootCmd.AddCommand(backupCmd)
}
