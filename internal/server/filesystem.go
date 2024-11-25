package server

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func InitRemotePath() (*cobra.Command, error) {
    return handleServerCmds()
}

func handleServerCmds() (*cobra.Command, error) {

	var rootCmd = &cobra.Command{
		Use:   "cls",
		Short: "Command-Line file sharing",
	}

	var remotePathCmd = &cobra.Command{
		Use:     "path",
		Short:   "Set filepath for received files",
		Long:    "Use path to set a custom path for where received files should go.",
		Example: "cls path <custom/path>",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{
			"-p",
			"-P",
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := os.MkdirAll(args[0], os.ModePerm); err != nil {
				log.Fatalf("Error while creating directory '%v': %v", args[0], err)
			}
		},
	}

    rootCmd.AddCommand(remotePathCmd)
    
    return rootCmd, nil
}
