package cmd

import (
	// "fmt"
	// "os"
	"github.com/spf13/cobra"
	"github.com/sorend/gokmp/pkg/backup"
)

var BackupDirectory string
var BackupNsid string

var backupCmd = &cobra.Command{
	Use: "backup",
	Short: "Backup from Flickr",
	Long: `Backup pictures from Flickr to local disk`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := backup.Run(BackupNsid, BackupDirectory); err != nil {
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
