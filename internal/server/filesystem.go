package server

import (
	"log"
	"os"
	"github.com/spf13/cobra"
)

// ExecuteRemotePath runs the server root command and returns a filepath string
func ExecuteRemotePath() string {
    var serverFilePath string

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
            fpath := args[0]
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				log.Fatalf("Error while creating directory '%v': %v", args[0], err)
			}
            serverFilePath = fpath
		},
	}

	rootCmd.AddCommand(remotePathCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
		// return nil, err
	}

    return serverFilePath

}
